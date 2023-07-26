package connection

import (
	"net"
	"sync"
)

type Client struct {
	ID   int64
	Conn net.Conn
	OS string
	Mu sync.Mutex
	IsActive bool
	Daddy string
}

type Command struct{
	CMD string
	Action string
	Response *Response
}

type Response struct{
	Data []any //! CHANGE LATER
	Result bool
}