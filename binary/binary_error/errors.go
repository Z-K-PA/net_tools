package binary_error

import "errors"

var (
	//序列化handler初始失败
	ErrInitHandler = errors.New("binary handler init error")

	//数据越界
	ErrOverflow = errors.New("binary handler data overflow")

	//字符串长度越界
	ErrStringOverflow = errors.New("binary handler string overflow")

	//数组越界
	ErrArrayOverflow = errors.New("binary handler array overflow")

	//反序号化的缓冲区为空
	ErrEmptyBuffer = errors.New("binary handler empty buffer")
)
