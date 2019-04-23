package base_connection

import (
	"net"
	"sync"
)

//解析接下来应该读取的字节数
type ParseNextReadNumFunc func(buf []byte) int

//基本连接
type BaseConnection struct {
	//连接
	conn net.Conn
	//先读的字节数
	readFirstNum int
	//解析接下来应该读取的字节数的函数
	nextReadNumFunc ParseNextReadNumFunc
	//是否关闭
	closed bool
	//锁
	sync.Mutex
}

//新建连接
func NewBaseConnection(conn net.Conn, readFirstNum int, nextReadNumFunc ParseNextReadNumFunc) *BaseConnection {
	return &BaseConnection{
		conn:conn,
		readFirstNum:readFirstNum,
		nextReadNumFunc:nextReadNumFunc,
	}
}
