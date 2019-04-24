package rpc

import (
	"github.com/pineal-niwan/busybox/buffer"
	"github.com/pineal-niwan/busybox/rpc/rpc_error"
	"github.com/pineal-niwan/busybox/util"
	"net"
	"sync"
	"time"
)

//由头部数据解析消息体大小的函数定义
// headBuf -- 头部字节流
// count -- 解析出来接下来应该读取的消息体大小
// err -- 解析错误
type ParseDataLenFunc func(headBuf []byte) (count int, err error)

//设置发送消息缓存的函数定义
// data -- 初始的缓冲区
// maxSize -- 消息内容的最大长度
// ---
// out -- 返回的缓冲区,如果data足够使用,则out就使用data,否则新开内存
// count -- 需要发送内容的大小
// err -- 错误信息
type SendHandler func(data []byte, maxSize int) (out []byte, count int, err error)

//设置处理消息缓存的函数定义
// data -- 缓冲区
// maxSize -- 消息内容的长度(含头部)
// ---
// out -- 返回的缓冲区,如果data足够使用,则out就使用data,否则新开内存,此缓存区会被发送回去
// count -- 需要发送内容的大小
// err -- 错误信息
type RecvHandler func(data []byte, inputSize int, maxSize int) (out []byte, count int, err error)

type Option struct {
	//头部字节流长度
	HeadSize int
	//初始化的bufferSize长度
	BufferSize int
	//buffer回收的限制长度
	//在有些RPC调用过程中,可能有会将buffer撑的比较大,
	//这里设定一个门槛,如果大于此门槛,则重新分配一次buffer,避免缓存浪费
	BufferRecycleSize int
	//最大的RPC消息体长度
	RPCMaxInputSize int
	//最大的RPC消息接收长度
	RPCMaxOutputSize int
	//解析出消息体内容大小的函数
	ParseDataSize ParseDataLenFunc
}

func (option *Option) Validate() error {
	if option.HeadSize == 0 ||
		option.BufferSize == 0 ||
		option.BufferRecycleSize == 0 ||
		option.RPCMaxInputSize == 0 ||
		option.RPCMaxOutputSize == 0 ||
		option.ParseDataSize == nil {
		return rpc_error.ErrInvalidRPCOption
	}
	if option.HeadSize > option.BufferSize {
		//连接收消息头的buffer大小都不够
		return rpc_error.ErrInvalidRPCOption
	}
	return rpc_error.ErrInvalidRPCOption
}

type Conn struct {
	conn net.Conn
	//option
	option Option
	//缓存
	buffer []byte
	//关闭
	closed bool
	sync.Mutex
}

//初始化
func (cli *Conn) Init(conn net.Conn, option Option) error {
	err := option.Validate()
	if err != nil {
		return err
	}
	cli.conn = conn
	cli.option = option
	cli.buffer = make([]byte, cli.option.BufferSize)
	return nil
}

//关闭
func (cli *Conn) Close() (err error) {
	cli.Lock()
	if cli.closed {
		cli.Unlock()
		return
	} else {
		cli.closed = true
		err = cli.conn.Close()
		cli.Unlock()
		return
	}
}

//发送和接收消息 -- 给客户端使用
//发送的消息缓存在buffer中，在发送完成后，又用缓存的buffer去接收
func (cli *Conn) SendAndRecv(
	sendHandler SendHandler,
	timeout time.Duration) (readSize int, err error) {

	//设置发送缓存区
	var sendSize int
	cli.buffer, sendSize, err = sendHandler(cli.buffer, cli.option.RPCMaxInputSize)
	if err != nil {
		return
	}

	if sendSize > cli.option.RPCMaxInputSize {
		//发送字节流太长
		err = rpc_error.ErrTooLongInputData
		return
	}

	//如果有超时,设置超时
	if timeout > 0 {
		err = cli.conn.SetDeadline(time.Now().Add(timeout))
		if err != nil {
			return
		}
	}

	//先发送数据
	err = util.NetSendBytes(cli.conn, cli.buffer[:sendSize])
	if err != nil {
		return
	}

	//读头部
	err = util.NetReadBytes(cli.conn, cli.buffer[:cli.option.HeadSize])
	if err != nil {
		return
	}
	readSize += cli.option.HeadSize

	//解析头部算出消息内容长度
	var dataSize int
	dataSize, err = cli.option.ParseDataSize(cli.buffer[:cli.option.HeadSize])
	if err != nil {
		return
	}

	if dataSize > cli.option.RPCMaxOutputSize {
		//接收到的消息体长度太长
		err = rpc_error.ErrTooLongOutputData
		return
	}
	if dataSize == 0 {
		//只有头没有内容
		err = rpc_error.ErrEmptyMsg
		return
	}
	//扩容buffer
	buffer.BytesExtends(cli.buffer, cli.option.HeadSize+dataSize)

	//读内容
	err = util.NetReadBytes(cli.conn, cli.buffer[cli.option.HeadSize:cli.option.HeadSize+dataSize])
	if err != nil {
		return
	}

	readSize += dataSize
	return
}

//接收和发送消息 -- 给服务端端使用
//接收的消息缓存在buffer中，在接收完成后，又用缓存的buffer去发送
//此连接复用,不需要设置超时
func (cli *Conn) RecvAndSend(recvHandler RecvHandler) (err error) {

	var readSize int
	//读头部
	err = util.NetReadBytes(cli.conn, cli.buffer[:cli.option.HeadSize])
	if err != nil {
		return
	}
	readSize += cli.option.HeadSize

	//解析头部算出消息内容长度
	var dataSize int
	dataSize, err = cli.option.ParseDataSize(cli.buffer[:cli.option.HeadSize])
	if err != nil {
		return
	}

	if dataSize > cli.option.RPCMaxInputSize {
		//接收到的消息体长度太长
		err = rpc_error.ErrTooLongInputData
		return
	}
	if dataSize == 0 {
		//只有头没有内容
		err = rpc_error.ErrEmptyMsg
		return
	}
	//扩容buffer
	buffer.BytesExtends(cli.buffer, cli.option.HeadSize+dataSize)

	//读内容
	err = util.NetReadBytes(cli.conn, cli.buffer[cli.option.HeadSize:cli.option.HeadSize+dataSize])
	if err != nil {
		return
	}
	readSize += dataSize

	//设置发送缓存区
	var sendSize int
	cli.buffer, sendSize, err = recvHandler(cli.buffer, readSize, cli.option.RPCMaxOutputSize)
	if err != nil {
		return
	}

	if sendSize > cli.option.RPCMaxOutputSize {
		//发送字节流太长
		err = rpc_error.ErrTooLongOutputData
		return
	}

	//发送数据
	err = util.NetSendBytes(cli.conn, cli.buffer[:sendSize])
	if err != nil {
		return
	}

	return
}

//获取缓存区
func (cli *Conn) GetData() []byte {
	return cli.buffer
}
