package util

import (
	"github.com/pineal-niwan/busybox/util/util_error"
	"io"
	"net"
)

//验证connection
func validate(conn net.Conn, buf []byte) (int, error) {
	if conn == nil {
		return 0, util_error.ErrEmptyConnection
	}

	bL := len(buf)
	if bL == 0 {
		return 0, util_error.ErrEmptyByteSlice
	}
	return bL, nil
}

//发送字节流
func NetSendBytes(conn net.Conn, buf []byte) error {
	size, err := validate(conn, buf)
	if err != nil {
		return err
	}
	writeSize := 0
	index := 0

	for {
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

//接收字节流
func NetReadBytes(conn net.Conn, buf []byte) error {
	_, err := validate(conn, buf)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(conn, buf)
	return err
}
