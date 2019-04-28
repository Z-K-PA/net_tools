package fast_rpc

import (
	"context"
	"fmt"
	"github.com/pineal-niwan/busybox/buffer"
	"github.com/pineal-niwan/busybox/util"
	"go.uber.org/zap"
	"sync"
)

type Cli struct {
	*CliOption
	connPool   *util.NetPool
	logger     *zap.Logger
	bufferPool sync.Pool
}

func (c *Cli) Call(ctx context.Context, inMsg IMsg) (outMsg IMsg, err error) {
	var size int

	conn, err := c.connPool.Get(ctx)
	if err != nil {
		return nil, err
	}

	buf := c.bufferPool.Get().([]byte)
	size, buf, err = inMsg.Marshal(buf, c.Option)
	if err != nil {
		return nil, err
	}

	if size > c.MaxMsgSize {
		return nil, fmt.Errorf("rpc client too long msg size:%+v", size)
	}

	err = util.NetSendBytes(conn, buf[:size])
	if err != nil {
		return nil, err
	}

	/***********************接收消息体***************/
	//接收消息体字节流
	var head MsgHead

	size = int(head.Size)
	if size == 0 || size > c.MaxMsgSize {
		c.logger.Error("client rpc size of outMsg error",
			zap.Int("msgSize", size))
		return
	}
	//如果buf不够，扩大
	buf = buffer.BytesExtends(buf, MsgHeadSize+size, 0)
	err = util.NetReadBytes(conn, buf[MsgHeadSize:MsgHeadSize+size])
	if err != nil {
		c.logger.Error("clint rpc receive content error",
			zap.Error(err))
		return
	}
	//解析消息内容
	outMsg, err = c.ParseMsg(head, buf[MsgHeadSize:MsgHeadSize+size])
	if err != nil {
		c.logger.Error("client parse content error",
			zap.Error(err))
		return
	}
	return
}
