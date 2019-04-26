package rpc

import (
	"context"
	"fmt"
	"github.com/pineal-niwan/busybox/rpc/rpc_error"
	"net"
	"sync"
)

//连接池对象
type ConnObj struct {
	*Conn
	p *Pool
}

//Return
func (c *ConnObj) Return() error {
	return c.p.put(c.Conn)
}

//连接池
type Pool struct {
	//conn chan
	connChan chan *Conn
	//连接器
	dialer *net.Dialer
	//服务地址
	address string
	//option
	option Option

	sync.RWMutex
}

//新建连接池
func NewPool(
	ctx context.Context,
	initCap int,
	dialer *net.Dialer,
	addr string,
	option Option) (*Pool, error) {

	if initCap <= 0 || dialer == nil {
		return nil, rpc_error.ErrInvalidPool
	}

	err := option.Validate()
	if err != nil {
		return nil, err
	}

	pool := &Pool{
		connChan: make(chan *Conn, initCap),
		dialer:   dialer,
		address:  addr,
		option:   option,
	}

	for i := 0; i < initCap; i++ {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("pool init err:%+v", err)
		}
		connNode := &Conn{}
		err = connNode.Init(conn, option)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("pool init err:%+v", err)
		}

		pool.connChan <- connNode
	}

	return pool, nil
}

//关闭
func (p *Pool) Close() error {
	var lastErr error

	p.Lock()
	connChan := p.connChan
	p.connChan = nil
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

//返回
func (p *Pool) put(conn *Conn) error {
	p.RLock()
	if p.connChan == nil {
		err := conn.Close()
		p.RUnlock()
		return err
	} else {
		p.connChan <- conn
		p.RUnlock()
		return nil
	}
}

//获取新建连接的参数
func (p *Pool) getFactory() (*net.Dialer, string, Option, chan *Conn) {
	p.RLock()
	dialer, addr, option, connChan := p.dialer, p.address, p.option, p.connChan
	p.RUnlock()
	return dialer, addr, option, connChan
}

//获取
func (p *Pool) Get(ctx context.Context) (*ConnObj, error) {
	dialer, addr, option, connChan := p.getFactory()
	if connChan == nil {
		return nil, rpc_error.ErrPoolClosed
	}

	select {
	case conn := <-connChan:
		if conn == nil {
			return nil, rpc_error.ErrPoolClosed
		} else {
			if conn.IsClosed() {
				//已经被关闭了的连接
				newConn, err := dialer.DialContext(ctx, "tcp", addr)
				if err != nil {
					//连接不上，保持chan大小不变，将原来的conn塞回去
					p.put(conn)
					return nil, err
				} else {
					newConnObj := &Conn{}
					err = newConnObj.Init(newConn, option)
					if err != nil {
						//连接不上，保持chan大小不变，将原来的conn塞回去
						p.put(conn)
						return nil, err
					} else {
						return &ConnObj{Conn: newConnObj, p: p}, nil
					}
				}
			} else {
				return &ConnObj{Conn: conn, p: p}, nil
			}
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
