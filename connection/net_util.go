package connection

import (
	"errors"
	"net"
)

var (
	ErrEmptyConnection = errors.New("empty connection")

	ErrEmptyByteSlice = errors.New("empty byte slice")
)

func validate(conn net.Conn, buf []byte) (int, error) {
	if conn == nil {
		return 0, ErrEmptyConnection
	}

	bL := len(buf)
	if bL == 0 {
		return 0, ErrEmptyByteSlice
	}
	return bL, nil
}

func SendBytes(conn net.Conn, buf []byte) error {
	size, err := validate(conn, buf)
	if err != nil {
		return err
	}
	writeSize := 0
	index := 0

	for{
		writeSize, err = conn.Write(buf[index:])
		if err != nil {
			return err
		}
		index += writeSize
		if index == size {
			return nil
		}
	}
}
