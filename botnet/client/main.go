package main

import (
	c "client/util/connection"
)

func main(){
	client := c.NewClient()
	client.ClientConnect()
}