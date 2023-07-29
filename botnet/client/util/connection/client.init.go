package connection

import (
	"net"
	"runtime"
	"strconv"
	"fmt"
	h "client/util/header"
)

func NewClient() *Client {
	return &Client{}
}

//*Sets the OS property to the runtime OS of the  client
func (c *Client) setOS(){
	c.OS = runtime.GOOS
}

//*Sets the ID property to input
func (c *Client) setID(id string){
	num, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		fmt.Println(h.E, "Error converting string to int64:", err)
		return
	}

	c.ID = num
}

//*Sets the Conn Propety to the param given
func (c *Client) setConn(conn net.Conn){
	c.Conn = conn
}