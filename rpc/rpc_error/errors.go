package rpc_error

import "errors"

var (
	ErrInvalidRPCOption = errors.New("invalid rpc option")

	ErrTooLongInputData = errors.New("too long input data")

	ErrTooLongOutputData = errors.New("too long output data")

	ErrEmptyMsg = errors.New("empty rpc content")

	ErrEmptyLogger = errors.New("empty logger")
)
