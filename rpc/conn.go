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
		return rpc_error.ErrInvalidOption
	}
	if option.HeadSize > option.BufferSize {
		//连接收消息头的buffer大小都不够
		return rpc_error.ErrInvalidOption
	}

	return nil
}

type Conn struct {
	//连接
	net.Conn
	//option
	option Option
	//缓存
	buffer []byte
	//关闭
	closed bool
	sync.RWMutex
}

//初始化
func (c *Conn) Init(conn net.Conn, option Option) error {
	if conn == nil {
		return rpc_error.ErrEmptyConnection
	}
	err := option.Validate()
	if err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			return closeErr
		}
		return err
	}
	c.Conn = conn
	c.option = option
	c.buffer = make([]byte, c.option.BufferSize)
	return nil
}

//关闭
func (c *Conn) Close() (err error) {
	c.Lock()
	if c.closed {
		c.Unlock()
		return
	} else {
		c.closed = true
		err = c.Conn.Close()
		c.Unlock()
		return
	}
}

//是否关闭
func (c *Conn) IsClosed() bool {
	c.RLock()
	closed := c.closed
	c.RUnlock()
	return closed
}

//发送和接收消息 -- 给客户端使用
//发送的消息缓存在buffer中，在发送完成后，又用缓存的buffer去接收
func (c *Conn) SendAndRecv(
	sendHandler SendHandler,
	deadline time.Time) (readSize int, err error) {

	//设置发送缓存区
	var sendSize int
	c.buffer, sendSize, err = sendHandler(c.buffer, c.option.RPCMaxInputSize)
	if err != nil {
		return
	}

	if sendSize > c.option.RPCMaxInputSize {
		//发送字节流太长
		err = rpc_error.ErrTooLongInputData
		return
	}

	//如果有超时,设置超时
	if !deadline.IsZero() {
		err = c.Conn.SetDeadline(deadline)
		if err != nil {
			return
		}
	}

	//先发送数据
	err = util.NetSendBytes(c.Conn, c.buffer[:sendSize])
	if err != nil {
		return
	}

	//读头部
	err = util.NetReadBytes(c.Conn, c.buffer[:c.option.HeadSize])
	if err != nil {
		return
	}
	readSize += c.option.HeadSize

	//解析头部算出消息内容长度
	var dataSize int
	dataSize, err = c.option.ParseDataSize(c.buffer[:c.option.HeadSize])
	if err != nil {
		return
	}

	if dataSize > c.option.RPCMaxOutputSize {
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
	buffer.BytesExtends(c.buffer, c.option.HeadSize+dataSize, 0)

	//读内容
	err = util.NetReadBytes(c.Conn, c.buffer[c.option.HeadSize:c.option.HeadSize+dataSize])
	if err != nil {
		return
	}

	readSize += dataSize
	return
}

//接收和发送消息 -- 给服务端端使用
//接收的消息缓存在buffer中，在接收完成后，又用缓存的buffer去发送
//此连接复用,不需要设置超时
func (c *Conn) RecvAndSend(recvHandler RecvHandler) (err error) {

	var readSize int
	//读头部
	err = util.NetReadBytes(c.Conn, c.buffer[:c.option.HeadSize])
	if err != nil {
		return
	}
	readSize += c.option.HeadSize

	//解析头部算出消息内容长度
	var dataSize int
	dataSize, err = c.option.ParseDataSize(c.buffer[:c.option.HeadSize])
	if err != nil {
		return
	}

	if dataSize > c.option.RPCMaxInputSize {
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
	buffer.BytesExtends(c.buffer, c.option.HeadSize+dataSize, 0)

	//读内容
	err = util.NetReadBytes(c.Conn, c.buffer[c.option.HeadSize:c.option.HeadSize+dataSize])
	if err != nil {
		return
	}
	readSize += dataSize

	//设置发送缓存区
	var sendSize int
	c.buffer, sendSize, err = recvHandler(c.buffer, readSize, c.option.RPCMaxOutputSize)
	if err != nil {
		return
	}

	if sendSize > c.option.RPCMaxOutputSize {
		//发送字节流太长
		err = rpc_error.ErrTooLongOutputData
		return
	}

	//发送数据
	err = util.NetSendBytes(c.Conn, c.buffer[:sendSize])
	if err != nil {
		return
	}

	return
}

//获取缓存区
func (c *Conn) GetData() []byte {
	return c.buffer
}

//Recycle buffer
func (c *Conn) Recycle() {
	if len(c.buffer) > c.option.BufferRecycleSize {
		c.buffer = make([]byte, c.option.BufferSize)
	}
}
