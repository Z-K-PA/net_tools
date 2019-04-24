package rpc

import (
	"github.com/pineal-niwan/busybox/rpc/rpc_error"
	"github.com/pineal-niwan/busybox/util"
	"github.com/pineal-niwan/busybox/util/util_error"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type Service struct {
	ln      net.Listener
	logger  *zap.Logger
	option  Option
	handler RecvHandler
	closed  bool
	sync.Mutex
}

//监听端口并服务
func (s *Service) ListenAndServer(address string, option Option, logger *zap.Logger, handler RecvHandler) (err error) {

	if logger == nil {
		return rpc_error.ErrEmptyLogger
	}

	err = option.Validate()
	if err != nil {
		return
	}

	s.Lock()
	s.ln, err = net.Listen("tcp", address)
	if err != nil {
		s.Unlock()
		return
	}
	s.Unlock()
	s.logger = logger

	go s.LoopListenServer()
	return
}

//循环监听端口
func (s *Service) LoopListenServer() {
	var conn net.Conn
	var accDelay time.Duration
	var accRetryCount int
	var err error

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
						accDelay = ACCEPT_DELAY
					} else {
						accDelay *= 2
					}
					if accDelay >= ACCEPT_MAX_DELAY {
						accDelay = ACCEPT_MAX_DELAY
					}
					time.Sleep(accDelay)
					accRetryCount++
					if accRetryCount >= ACCEPT_MAX_RETRY {
						//超过重试次数
						s.logger.Error(
							"RPC Service超过重试次数",
							zap.Int("retry", accRetryCount))
						return
					}
					//可以继续
					continue
				} else {
					//不是临时错误
					s.logger.Error("RPC Service不是临时错误", zap.Error(err))
					return
				}
			} else {
				//不是网络错误
				s.logger.Error("RPC 不是网络错误, err", zap.Error(err))
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

//关闭监听端口
func (s *Service) CloseListener() (err error) {
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

//处理连接
func (s *Service) HandleConnection(conn net.Conn) {
	cli := Conn{}
	err := cli.Init(conn, s.option)
	if err != nil {
		s.logger.Error("初始化进入的连接失败",
			zap.Error(err))
		return
	}

	defer func() {
		//panic后防止整个server被panic
		panicErr := util.Recover(recover())
		if panicErr != nil {
			pErr := util_error.NewPanicError()
			s.logger.Error("service panic",
				zap.Error(pErr))
			s.logger.Error("service panic error",
				zap.Error(panicErr.Err))
			s.logger.Error("service panic stack:",
				zap.String("stack", string(panicErr.Stack())))
		}
		//关闭连接
		err = cli.Close()
		if err != nil {
			s.logger.Error("service close connection", zap.Error(err))
		}
	}()
	for {
		//不停地收发处理消息
		err = cli.RecvAndSend(s.handler)
		if err != nil {
			s.logger.Error("service handler", zap.Error(err))
		}
	}
}
