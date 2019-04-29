package main

import (
	"errors"
	"github.com/pineal-niwan/busybox/util"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
)

var (
	ErrNoAddress = errors.New("not specific address")
)

//server_run
func serverRun(c *cli.Context, logger *zap.Logger) error {
	//参数检查
	address := c.String("address")
	pprofAddress := c.String("pprofAddress")

	if address == "" || pprofAddress == "" {
		return ErrNoAddress
	}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	//rpc notify chan
	rpcNotify := make(chan struct{})
	//启动rpc服务
	service, err := initServiceHandler(ln, logger)
	if err != nil {
		//关闭连接
		ln.Close()
		return err
	}
	go service.LoopHandle(rpcNotify)

	//pprof notify chan
	pprofNotify := make(chan error)
	go util.PprofServerStart(pprofAddress, pprofNotify)

	//等待信号
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, os.Kill)

	select {
	case <-rpcNotify:
		logger.Error("rpc listener exit")
	case pprofErr := <-pprofNotify:
		logger.Error("pprof exit", zap.Error(pprofErr))
	case <-kill:
		logger.Error("server killed by signal")
	}

	return nil
}
