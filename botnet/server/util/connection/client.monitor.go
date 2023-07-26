package connection

import (
	"fmt"
	h "server/util/header"
)



/*
* Go Routine that Checks if a client is disconnected
* if value in chan remove client from list
*/
func (s *Server) ClientStatus() {
	for {
		disconnectedClient := <-s.DisconnectedClients
		fmt.Println(h.E, "Client ID", disconnectedClient.ID, "Disconnected")
		s.removeClientByID(disconnectedClient.ID)
		s.resetActive(disconnectedClient.ID)
	}
}

func (s *Server) resetActive(clientID int64){
	s.mu.Lock()
	defer s.mu.Unlock()

	itemMap := make(map[int64]bool)
	for _, item := range s.Clients{
		itemMap[item.ID] = true
	}
	
	
	if !itemMap[clientID]{
		fmt.Println(h.I, "Active Client has been reset to 0 (default)")
		s.ActiveClient = 0
	}
}

/*
*removes client from list
*/
func (s *Server) removeClientByID(clientID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var updatedClients []*Client

	for _, client := range s.Clients {
		if client.ID != clientID {
			updatedClients = append(updatedClients, client)
		}
	}

	s.Clients = updatedClients
}

/*
* Go Routine to update channel when client disconnects
*/
func (s *Server) HandleClientDisconnection(client *Client) {
	buffer := make([]byte, 1)
	_, err := client.Conn.Read(buffer)
	if err != nil {
		s.DisconnectedClients <- client
	}
}

