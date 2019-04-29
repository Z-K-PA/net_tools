package fast_rpc

import (
	"github.com/pineal-niwan/busybox/binary"
	"time"
)

type Option struct {
	//序列化选项
	*binary.Option

	//连接accept出临时错误时的延时
	AcceptDelay time.Duration
	//连接accept出临时错误时的最大延时
	AcceptMaxDelay time.Duration
	//连接accept出临时错误时的最大重试次数
	AcceptMaxRetry int
	//为每个连接初始分配的字节数
	BufferSize int
	//最大的消息体长度
	MaxMsgSize int
	//每个连接的buffer回收门槛
	BufferRecycleSize int
}

func (option *Option) Validate() error {
	if !option.Option.Validate() {
		return ErrInvalidOption
	}

	if option.AcceptDelay == 0 ||
		option.AcceptMaxDelay == 0 ||
		option.AcceptMaxRetry == 0 ||
		option.BufferSize == 0 ||
		option.MaxMsgSize == 0 ||
		option.BufferRecycleSize == 0 {
		return ErrInvalidOption
	}
	return nil
}

type CliOption struct {
	//序列化选项
	*binary.Option

	//为每个连接初始分配的字节数
	BufferSize int
	//最大的消息体长度
	MaxMsgSize int
	//每个连接的buffer回收门槛
	BufferRecycleSize int
	//退火时间
	RetreatTime time.Duration
}

func (cliOption *CliOption) Validate() error {
	if !cliOption.Option.Validate() {
		return ErrInvalidOption
	}

	if cliOption.BufferSize == 0 ||
		cliOption.MaxMsgSize == 0 ||
		cliOption.BufferRecycleSize == 0 {
		return ErrInvalidOption
	}
	return nil
}
