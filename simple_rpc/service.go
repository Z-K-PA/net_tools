package simple_rpc

import (
	"errors"
	"github.com/pineal-niwan/busybox/util"
	"github.com/pineal-niwan/busybox/util/util_error"
	"go.uber.org/zap"
	"log"
	"net"
	"sync"
	"time"
)

const (
	//连接accept出临时错误时的延时
	AcceptDelay    = time.Microsecond * 5
	//连接accept出临时错误时的最大延时
	AcceptMaxDelay = time.Millisecond * 200
	//连接accept出临时错误时的最大重试次数
	AcceptMaxRetry = 1000

	//为每个连接初始分配的字节数
	BufferSize = 1024
	//最大的消息体长度
	MaxMsgSize = 1024*1024
	//每个连接的buffer回收门槛
	BufferRecycleSize = 1024*1024
)

var (
	//空logger错
	ErrEmptyLogger = errors.New("empty logger")
	//buffer不够
	ErrNotEnoughBuffer = errors.New("not enough buffer")
)

type Service struct {
	ln     net.Listener
	logger *zap.Logger
	closed bool
	sync.Mutex
}

//监听端口并服务
func (s *Service) ListenAndServer(address string, logger *zap.Logger) (err error) {

	if logger == nil {
		return ErrEmptyLogger
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
						accDelay = AcceptDelay
					} else {
						accDelay *= 2
					}
					if accDelay >= AcceptMaxDelay {
						accDelay = AcceptMaxDelay
					}
					time.Sleep(accDelay)
					accRetryCount++
					if accRetryCount >= AcceptMaxRetry {
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
	var msg IMsg
	buf := make([]byte, BufferSize)

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
		closeErr := conn.Close()
		if closeErr != nil {
			s.logger.Error("service close connection", zap.Error(closeErr))
		}
	}()
	for {
		//不停地收发处理消息
		if len(buf) < MsgHeadSize {
			s.logger.Error("no enough buffer")
			s.logger.Sync()
			log.Fatal("not enough buffer in service")
		}
		err := util.NetReadBytes(conn, buf[:MsgHeadSize])
		if err != nil {
			s.logger.Error("service receive head error",
				zap.Error(err))
			return
		}



		if err != nil {
			s.logger.Error("service handler", zap.Error(err))
		}
	}
}
