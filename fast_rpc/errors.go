package fast_rpc

import "errors"

var (
	//未知错误-逻辑上不应该出现
	ErrUnknown       = errors.New("unknown error")
	//option参数不正确
	ErrInvalidOption = errors.New("invalid option")
	//没有设置消息解析器
	ErrBadMsgParser  = errors.New("bad msg parser")
	//没有设置消息处理器
	ErrBadMsgHandler = errors.New("bad msg handler")
	//接收的消息不是自己期望的
	ErrNotExpectMsg = errors.New("not expect message")
)
