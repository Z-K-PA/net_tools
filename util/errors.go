package util

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrEmptyConnection = errors.New("empty connection")

	ErrEmptyByteSlice = errors.New("empty byte slice")
)

func NewPanicError() error {
	var ok bool
	var fileName string
	var line int

	_, fileName, line, ok = runtime.Caller(2)
	if !ok {
		fileName, line = "???", 0
	}

	return fmt.Errorf("panic error in file %s : line: %d", fileName, line)
}
