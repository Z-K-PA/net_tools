package fast_rpc

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/pineal-niwan/busybox/binary"
	"github.com/pineal-niwan/busybox/buffer"
	"github.com/pineal-niwan/busybox/util"
	"go.uber.org/zap"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrBadMsgParser  = errors.New("bad msg parser")
	ErrBadMsgHandler = errors.New("bad msg parser")
)

//消息解析函数
type MsgParseHandler func(msgData []byte, option *binary.Option) (IMsg, error)

//消息处理函数
type MsgHandler func(inMsg IMsg) (outMsg IMsg, err error)

//服务定义
type Service struct {
	//监听端口
	ln net.Listener
	//日志
	logger *zap.Logger
	//参数
	option *Option
	//消息解析
	msgParseHash map[uint32]MsgParseHandler
	//消息处理
	msgHandlerHash map[uint32]MsgHandler

	closed bool
	sync.Mutex
}

//初始化
func (s *Service) Init(
	ln net.Listener,
	logger *zap.Logger,
	option *Option,
	msgParseHash map[uint32]MsgParseHandler) {
	s.Lock()
	s.ln = ln
	s.logger = logger
	s.option = option
	s.msgParseHash = msgParseHash
	s.msgHandlerHash = make(map[uint32]MsgHandler)
	s.Unlock()
}

//添加消息处理
func (s *Service) AddMsgHandler(msg IMsg, handler MsgHandler) {
	s.msgHandlerHash[msg.GetCode()] = handler
}

//关闭监听端口
func (s *Service) Close() (err error) {
	s.Lock()
	//已经关闭过了
	if s.closed {
		s.Unlock()
		return
	}
	s.closed = true
	//关闭
	if s.ln != nil {
		err = s.ln.Close()
	}
	s.Unlock()
	return
}

//循环处理消息
func (s *Service) LoopHandle(exitNotify chan<- struct{}) {
	var conn net.Conn
	var accDelay time.Duration
	var accRetryCount int
	var err error

	defer func() {
		exitNotify <- struct{}{}
	}()

	for {
		//监听socket
		conn, err = s.ln.Accept()
		if err != nil {
			ne, ok := err.(net.Error)
			if ok {
				//是网络错误
				if ne.Temporary() {
					//是临时错误，可以修复
					if accDelay <= 0 {
						accDelay = s.option.AcceptDelay
					} else {
						accDelay *= 2
					}
					if accDelay >= s.option.AcceptMaxDelay {
						accDelay = s.option.AcceptMaxDelay
					}
					time.Sleep(accDelay)
					accRetryCount++
					if accRetryCount >= s.option.AcceptMaxRetry {
						//超过重试次数
						s.logger.Error(
							"Service accept 超过重试次数",
							zap.Int("retry", accRetryCount))
						return
					}
					//可以继续
					continue
				} else {
					//不是临时错误
					s.logger.Error("Service accept 不是临时错误", zap.Error(err))
					return
				}
			} else {
				//不是网络错误
				s.logger.Error("Service accept 不是网络错误, err", zap.Error(err))
				return
			}
		}
		//没有错误，重置重试的变量
		accRetryCount = 0
		accDelay = 0

		//新加入连接进行处理
		go s.HandleConnection(conn)
	}
}

//处理连接
func (s *Service) HandleConnection(conn net.Conn) {
	var inMsg, outMsg IMsg
	var err error
	var head MsgHead
	var size int

	buf := make([]byte, s.option.BufferSize)

	defer func() {
		//panic后防止整个server被panic
		panicErr := util.Recover(recover())
		if panicErr != nil {
			pErr := util.NewPanicError()
			s.logger.Error("service panic",
				zap.Error(pErr))
			s.logger.Error("service panic error",
				zap.Error(panicErr.Err))
			s.logger.Error("service panic stack:",
				zap.String("stack", string(panicErr.Stack())))
		}
		//关闭连接
		closeErr := conn.Close()
		if closeErr != nil {
			s.logger.Error("service close connection", zap.Error(closeErr))
		}
	}()

	for {
		if len(buf) < MsgHeadSize {
			//缓存不够，退出
			s.logger.Error("no enough buffer")
			s.logger.Sync()
			log.Fatal("not enough buffer in service")
		}

		/***********************接收消息头***************/
		//接收消息头字节流
		err = util.NetReadBytes(conn, buf[:MsgHeadSize])
		if err != nil {
			s.logger.Error("service receive head error",
				zap.Error(err))
			return
		}
		//解析消息头
		head, err = UnmarshalMsgHead(buf, s.option.Option)
		if err != nil {
			s.logger.Error("service parse head error",
				zap.Error(err))
			return
		}

		/***********************接收消息体***************/
		//接收消息体字节流
		size = int(head.Size)
		if size == 0 || size > s.option.MaxMsgSize {
			s.logger.Error("service size of inMsg error",
				zap.Int("msgSize", size))
			return
		}
		//如果buf不够，扩大
		buf = buffer.BytesExtends(buf, MsgHeadSize+size, 0)
		err = util.NetReadBytes(conn, buf[MsgHeadSize:MsgHeadSize+size])
		if err != nil {
			s.logger.Error("service receive content error",
				zap.Error(err))
			return
		}
		//解析消息内容
		inMsg, err = s.ParseMsg(head, buf)
		if err != nil {
			s.logger.Error("service parse content error",
				zap.Error(err))
			return
		}

		/***********************处理消息****************/
		outMsg, err = s.HandleMsg(inMsg)
		if err != nil {
			s.logger.Error("service handle msg error",
				zap.Error(err))
			return
		}

		/***********************返回结果消息****************/
		buf, size, err = outMsg.Marshal(buf)
		if err != nil {
			s.logger.Error("service marshal out msg error",
				zap.Error(err))
			return
		}
		//发送字节流
		err = util.NetSendBytes(conn, buf[:size])
		if err != nil {
			s.logger.Error("service send out msg error",
				zap.Error(err))
			return
		}
	}

}

//解析消息
func (s *Service) ParseMsg(head MsgHead, buf []byte) (IMsg, error) {
	if s.msgParseHash == nil {
		return nil, ErrBadMsgParser
	}

	parseHandler, ok := s.msgParseHash[head.GetCode()]
	if !ok || parseHandler == nil {
		err := fmt.Errorf("bad msg parser cmd:%+v, version:%+v", head.Cmd, head.Version)
		return nil, err
	}
	return parseHandler(buf, s.option.Option)
}

//处理消息
func (s *Service) HandleMsg(inMsg IMsg) (IMsg, error) {
	if s.msgHandlerHash == nil {
		return nil, ErrBadMsgHandler
	}

	msgHandler, ok := s.msgHandlerHash[inMsg.GetCode()]
	if !ok || msgHandler == nil {
		err := fmt.Errorf("bad msg handler cmd:%+v, version:%+v", inMsg.GetCmd(), inMsg.GetVersion())
		return nil, err
	}
	return msgHandler(inMsg)
}
