package main

import (
	"github.com/pineal-niwan/busybox/fast_rpc"
	"go.uber.org/zap"
	"net"
)

func initServiceHandler(ln net.Listener, logger *zap.Logger) (*fast_rpc.Service, error) {

	var msgParseHandler map[uint32]fast_rpc.MsgParseHandler

	service := &fast_rpc.Service{}
	option := &fast_rpc.Option{
		//Add you init code here
	}

	err := option.Validate()
	if err != nil{
		return nil, err
	}

	// uncomment this code to init your msg parser
	//msgParseHandler = xxx

	service.Init(ln, logger, option, msgParseHandler)

	//add your service handler here
	//service.AddMsgHandler()

	return service, err
}
