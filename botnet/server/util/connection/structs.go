package connection

import (
	"net"
	"sync"
)

type Server struct {
	ActiveConns         int
	Clients             []*Client 
	ActiveClient        int64     // 0 = none & # = id selected currently
	ClientGroup         []*Client
	DisconnectedClients chan *Client 
	mu                  sync.Mutex // Mutex to manage concurrent access to s.Clients
}

type Client struct {
	ID      int64
	Conn net.Conn
	Status bool
}