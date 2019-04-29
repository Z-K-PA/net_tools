package fast_rpc

import (
	"context"
	"fmt"
	"github.com/pineal-niwan/busybox/buffer"
	"github.com/pineal-niwan/busybox/util"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type Cli struct {
	//参数
	*CliOption
	//日志
	logger *zap.Logger
	//连接池
	connPool *util.NetPool
	//缓冲池
	bufferPool *sync.Pool
	//消息解析
	msgParseHash map[uint32]MsgParseHandler
}

func NewCli(
	ctx context.Context,
	address string,
	poolSize int,
	option *CliOption,
	msgParseHash map[uint32]MsgParseHandler) (*Cli, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	netDialer := &net.Dialer{
		KeepAlive: 5 * time.Minute, //5分钟
	}
	connPool, err := util.NewPool(ctx, poolSize, netDialer, address)
	if err != nil {
		return nil, err
	}

	bufferPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, option.BufferSize)
		},
	}

	cli := &Cli{
		CliOption:    option,
		logger:       logger,
		connPool:     connPool,
		bufferPool:   bufferPool,
		msgParseHash: msgParseHash,
	}
	return cli, nil
}

//多次调用
func (cli *Cli) CallWithRetry(ctx context.Context, inMsg IMsg, retryTimes int) (IMsg, error) {
	var conn *util.Conn
	var callRet *_CallRet
	var err error

	//先去连接池拿连接
	conn, err = cli.connPool.Get(ctx)
	if err != nil {
		//拿不到连接，直接退出
		return nil, err
	}

	//拿到连接后才开始分配缓冲区
	buf := cli.bufferPool.Get().([]byte)
	defer func() {
		//捕获panic
		panicErr := util.Recover(recover())
		if panicErr != nil {
			pErr := util.NewPanicError()
			cli.logger.Error("client rpc panic",
				zap.Error(pErr))
			cli.logger.Error("client rpc panic error",
				zap.Error(panicErr.Err))
			cli.logger.Error("client rpc panic stack:",
				zap.String("stack", string(panicErr.Stack())))
		}
		//归还可复用的缓冲区
		if len(buf) <= cli.BufferRecycleSize {
			cli.bufferPool.Put(buf)
		}
		//归还连接
		conn.Close()
	}()

	callRet = cli.callWithConn(ctx, conn, inMsg, buf)
	if callRet.err == nil {
		//一次调用就成功了
		return callRet.msg, callRet.err
	}

	if !callRet.needResetConn {
		//逻辑错误 -- 重连也是枉然
		return nil, callRet.err
	}

	for i := 0; i < retryTimes; i++ {
		if cli.RetreatTime > 0 {
			//暂时等待
			time.Sleep(cli.RetreatTime * time.Duration(i+1))
		}
		err = conn.Renew(ctx)
		if err != nil {
			continue
		}
		//走到这里表明重连成功
		//继续调用
		callRet = cli.callWithConn(ctx, conn, inMsg, buf)
		if callRet.err == nil {
			return callRet.msg, callRet.err
		} else {
			//继续调用失败
			if callRet.needResetConn {
				//可以继续重试
				err = callRet.err
				continue
			} else {
				return nil, callRet.err
			}
		}
	}

	if err != nil {
		return nil, err
	} else {
		return nil, ErrUnknown
	}
}

type _CallRet struct {
	//返回的消息
	msg IMsg
	//错误
	err error
	//是否需要重置连接
	needResetConn bool
	//复用的缓存
	buf []byte
}

func (ret *_CallRet) set(msg IMsg, err error, needResetConn bool, buf []byte) *_CallRet {
	ret.msg = msg
	ret.err = err
	ret.needResetConn = needResetConn
	ret.buf = buf
	return ret
}

//调用RPC - 带conn
//返回值 (IMsg -- 返回的消息 error-错误 bool-是否需要重置连接)
func (cli *Cli) callWithConn(ctx context.Context, conn net.Conn, inMsg IMsg, buf []byte) *_CallRet {
	var err error
	var size int

	callRet := &_CallRet{}

	/***********************序列化消息***************/
	size, buf, err = inMsg.Marshal(buf, cli.Option)
	if err != nil {
		callRet.msg = nil
		callRet.err = err
		return callRet.set(nil, err, false, buf)
	}

	if size > cli.MaxMsgSize {
		return callRet.set(
			nil,
			fmt.Errorf("rpc client too long msg size:%+v", size),
			false,
			buf)
	}

	/***********************发送消息体***************/
	deadline, ok := ctx.Deadline()
	if ok {
		err = conn.SetDeadline(deadline)
		if err != nil {
			return callRet.set(nil, err, true, buf)
		}
	}

	err = util.NetSendBytes(conn, buf[:size])
	if err != nil {
		return callRet.set(nil, err, true, buf)
	}

	/***********************接收消息头***************/
	var head MsgHead
	//接收消息头字节流
	err = util.NetReadBytes(conn, buf[:MsgHeadSize])
	if err != nil {
		cli.logger.Error("rpc client receive head error",
			zap.Error(err))
		return callRet.set(nil, err, true, buf)
	}
	//解析消息头
	head, err = UnmarshalMsgHead(buf, cli.Option)
	if err != nil {
		cli.logger.Error("rpc client parse head error",
			zap.Error(err))
		return callRet.set(nil, err, true, buf)
	}

	/***********************接收消息体***************/
	//检查消息体大小
	size = int(head.Size)
	if size == 0 || size > cli.MaxMsgSize {
		cli.logger.Error("client rpc size of outMsg error",
			zap.Int("msgSize", size))
		return callRet.set(
			nil,
			fmt.Errorf("msg size out of range :%+v", size),
			true,
			buf)
	}
	//如果buf不够，扩大
	buf = buffer.BytesExtends(buf, MsgHeadSize+size, 0)
	//接收消息体内容字节流
	err = util.NetReadBytes(conn, buf[MsgHeadSize:MsgHeadSize+size])
	if err != nil {
		cli.logger.Error("clint rpc receive content error",
			zap.Error(err))
		return callRet.set(nil, err, true, buf)
	}

	/***********************解析返回消息体*************/
	//解析消息内容
	outMsg, err := cli.ParseMsg(head, buf[MsgHeadSize:MsgHeadSize+size])
	if err != nil {
		cli.logger.Error("client parse content error",
			zap.Error(err))
		return callRet.set(nil, err, false, buf)
	}
	return callRet.set(outMsg, nil, false, buf)
}

//解析消息
func (cli *Cli) ParseMsg(head MsgHead, buf []byte) (IMsg, error) {
	if cli.msgParseHash == nil {
		return nil, ErrBadMsgParser
	}

	parseHandler, ok := cli.msgParseHash[head.GetCode()]
	if !ok || parseHandler == nil {
		err := fmt.Errorf("bad msg parser cmd:%+v, version:%+v", head.Cmd, head.Version)
		return nil, err
	}
	return parseHandler(buf, cli.Option)
}
