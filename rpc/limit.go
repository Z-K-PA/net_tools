package rpc

import "time"

const (
	ACCEPT_DELAY     = time.Microsecond * 5
	ACCEPT_MAX_DELAY = time.Millisecond * 200
	ACCEPT_MAX_RETRY = 1000
)
