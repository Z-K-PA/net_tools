package rpc

import "sync"

type Pool struct {
	connChan chan *ReusableConn
	sync.RWMutex
}

//可重复使用Conn
type ReusableConn struct {
	//连接
	*Conn
	//pool
	pool *Pool
}
