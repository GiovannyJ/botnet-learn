package connection

import (
	"net"
	h "client/util/header"
	"log"
)

func HandleClient(command string, conn net.Conn){
	switch {
	case command == "count":
		break
	case command == "detect os":
		break
	case command == "open application":
		break
	case command == "send file":
		break
	case command == "search file":
		break
	case command == "bye":
		break
	default:
		if _, err := conn.Write([]byte(h.E + " Unknown command: " + command + "\n")); err != nil {
			log.Println("Error:", err)
		}
	}
}