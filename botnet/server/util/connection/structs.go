package connection

import (
	"net"
	"sync"
)

type Server struct {
	Activeconns         int
	Clients             []*Client 
	ActiveClient        int64     // 0 = none & # = id selected currently
	ClientGroup         []int64
	DisconnectedClients chan *Client 
	mu                  sync.Mutex // Mutex to manage concurrent access to s.Clients
}

type Client struct {
	ID      int64
	Conn net.Conn
	Status bool
}