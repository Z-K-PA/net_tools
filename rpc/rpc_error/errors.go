package rpc_error

import "errors"

var (
	ErrInvalidOption = errors.New("invalid option")

	ErrEmptyConnection = errors.New("invalid empty connection")

	ErrTooLongInputData = errors.New("too long input data")

	ErrTooLongOutputData = errors.New("too long output data")

	ErrEmptyMsg = errors.New("empty rpc content")

	ErrEmptyLogger = errors.New("empty logger")

	ErrInvalidPool = errors.New("pool invalid settings")

	ErrPoolInitErr = errors.New("pool init error")

	ErrPoolClosed = errors.New("pool is closed")
)
