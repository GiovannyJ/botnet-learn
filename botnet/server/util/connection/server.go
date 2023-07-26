package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	conf "server/util/configuration"
	h "server/util/header"
	"strings"
	// "time"
)

/*
*Start TCP connection
* This file will be where the server hosts it main connection
* hosts on localhost:PORT
*/

//*Returns a new server instance, DisconnectedClients property initialized
func NewServer() *Server {
	return &Server{
		DisconnectedClients: make(chan *Client),
	}
}

//*current server instance used to track all data
var idGen = NewUniqueIDGenerator()


/*
*Starting the server Instance
*/
func (s *Server)ServerStart(mode string){
	HOST := conf.EnvVar("HOST")
	PORT := conf.EnvVar("PORT")
	TYPE := conf.EnvVar("TYPE")
	
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	
	if err != nil{
		fmt.Println(h.E, err)
		os.Exit(1)
	}

	fmt.Println(h.K,"Server started")
	fmt.Println(h.I, "Session will begin when client connects")
	fmt.Println(h.Line)
	
	//*GO ROUTINE TO CHECK IF ANY CLIENTS HAVE DISCONNECTED
	go s.ClientStatus()
	
	//*ITER LOOP OF SERVER
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println(h.K,"New Client Connected!")

		newClient := &Client{
			ID:      idGen.NextID(),
			Conn: conn,
		}
		
		s.mu.Lock()
		s.Clients = append(s.Clients, newClient)
		s.mu.Unlock()

		//*Interactive mode = send commands to clients
		if mode == "i"{
			fmt.Println(h.K, "interactive mode")
			go s.handleConnection(conn, newClient)
			go s.sendCommands(conn)
		}else{
			go s.handleConnection(conn, newClient)
		}
	}
}

/*
*Monitor and logs clients
*/
func (s *Server) handleConnection(conn net.Conn, client *Client) {
	defer client.Conn.Close()
	
	s.Activeconns++

	defer func() {
		s.Activeconns--
		s.DisconnectedClients <- client
	}()

	for {
		
		buffer := make([]byte, 1024)
		n, err := client.Conn.Read(buffer)
		
		if err != nil {
			fmt.Println(h.E, "Client Disconnected", err)
			break
		}

		contents := strings.TrimSpace(string(buffer[:n]))
		fmt.Println(h.I, "Client says:", contents)
		// HandleClient(s, contents, conn)	
		s.HandleClientDisconnection(client)
		
	}
}


/*
*Send commands to client
*/
func (s *Server) sendCommands(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		handleServer(s, input, conn)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading standard input:", err)
	}
}