package util

import (
	"context"
	"errors"
	"net"
	"sync"
)

var (
	ErrInvalidPool = errors.New("pool invalid settings")
	ErrPoolClosed  = errors.New("pool is closed")
)

//封装的net.Conn
type Conn struct {
	net.Conn
	p *NetPool
}

//Close
func (c *Conn) Close() error {
	return c.p.put(c.Conn)
}

//Renew
func (c *Conn) Renew(ctx context.Context) error {
	dialer, addr, connList := c.p.getFactory()
	if connList == nil {
		return ErrPoolClosed
	}
	newConn, err := dialer.DialContext(ctx, `tcp`, addr)
	if err != nil {
		return err
	}
	oldConn := c.Conn
	c.Conn = newConn
	oldConn.Close()
	return nil
}

//连接池
type NetPool struct {
	//conn chan list
	connList chan net.Conn
	//连接器
	dialer *net.Dialer
	//服务地址
	address string
	//锁
	sync.RWMutex
}

func NewPool(ctx context.Context, size int, dialer *net.Dialer, address string) (*NetPool, error) {
	if size <= 0 || dialer == nil {
		return nil, ErrInvalidPool
	}

	pool := &NetPool{
		connList: make(chan net.Conn, size),
		dialer:   dialer,
		address:  address,
	}

	for i := 0; i < size; i++ {
		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err != nil {
			pool.Close()
			return nil, err
		}

		pool.connList <- conn
	}
	return pool, nil
}

func (p *NetPool) Close() error {
	var lastErr error

	p.Lock()
	connChan := p.connList
	p.connList = nil
	p.Unlock()

	close(connChan)
	for conn := range connChan {
		err := conn.Close()
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (p *NetPool) put(conn net.Conn) error {
	p.RLock()
	if p.connList == nil {
		err := conn.Close()
		p.RUnlock()
		return err
	} else {
		p.connList <- conn
		p.RUnlock()
		return nil
	}
}

func (p *NetPool) getFactory() (*net.Dialer, string, chan net.Conn) {
	p.RLock()
	dialer, addr, connList := p.dialer, p.address, p.connList
	p.RUnlock()
	return dialer, addr, connList
}

func (p *NetPool) Get(ctx context.Context) (net.Conn, error) {
	_, _, connList := p.getFactory()

	if connList == nil {
		return nil, ErrPoolClosed
	}

	select {
	case conn := <-connList:
		if conn == nil {
			return nil, ErrPoolClosed
		} else {
			return &Conn{
				Conn: conn,
				p:    p,
			}, nil
		}
	}
}
