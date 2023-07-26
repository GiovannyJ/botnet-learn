package connection

import (
	"fmt"
	"net"
	h "client/util/header"
	// "strconv"
	// "strings"
	"os"
	conf "client/util/config"
	// "io"
)

func NewClient() *Client{
	return &Client{}
}


func (c *Client) ClientConnect(){
	HOST := conf.EnvVar("HOST")
	PORT := conf.EnvVar("PORT")
	TYPE := conf.EnvVar("TYPE")

    // Establish a TCP connection to the server
    conn, err := net.Dial(TYPE, HOST+":"+PORT)
    if err != nil {
        fmt.Println(h.E, "Error connecting to server:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println(h.K, "Connected to server")

    for {
        // You can now send and receive data with the server using the "conn" variable.
        // For example, sending data:
        // message := "Hello, server!"
        // _, err = conn.Write([]byte(message))
        // if err != nil {
        //     fmt.Println("Error sending data:", err)
        //     break
        // }

        // Reading response from the server
        buffer := make([]byte, 1024)
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(h.E, "Error reading data:", err)
            break
        }

        // Convert the response to a string and print it
        command := string(buffer[:bytesRead])
        c.handleCommand(command, conn)
    }
	
}

func (c *Client) handleCommand(command string, conn net.Conn){
	fmt.Println(h.K, "Server response:", command)
}