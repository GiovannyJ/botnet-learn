package connection

import (
	"fmt"
	"io"
	// "net"
	"os"
	"strconv"
	"strings"
	h "server/util/header"
)

const BUFFERSIZE = 1024

//*Send File to client
// func (s *Server) SendFile(file string, client *Client) (bool, error){
// 	defer client.Conn.Close()
// 	f, err := os.OpenFile(file, 0, O_)
	
// 	return false, net.ErrClosed
// }

func (s *Server) SendFileToClient(f string, client *Client) error{
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := FillString(fileInfo.Name(), 64)
	fmt.Println(h.I, "Sending filename and filesize")
	
	client.Conn.Write([]byte(fileSize))
	client.Conn.Write([]byte(fileName))

	sendBuffer := make([]byte, BUFFERSIZE)

	fmt.Println(h.I, "Sending File to: ", client.ID)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		client.Conn.Write(sendBuffer)
	}

	fmt.Println(h.K, "File has been sent to:", client.ID)
	return nil
}







func (s *Server) RecvFileFromClient(client *Client) error{
	// connection, err := net.Dial("tcp", "localhost:27001")
	// if err != nil {
	// 	panic(err)
	// }
	// defer connection.Close()
	// fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	
	client.Conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
	client.Conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	
	newFile, err := os.Create(fileName)
	
	if err != nil {
		// panic(err)
		return err
	}
	defer newFile.Close()
	var receivedBytes int64
	
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, client.Conn, (fileSize - receivedBytes))
			client.Conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, client.Conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println(h.K, "Received file completely!")
	return nil
}