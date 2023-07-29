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
}

type Command struct{
	CMD string
	Action string
	FileInfo File
	Response *Response
}

type Response struct{
	Data []byte
	Result bool
}

type File struct{
	Name string
	Size int64
}