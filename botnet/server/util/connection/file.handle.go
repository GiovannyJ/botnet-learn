package connection

import (
	"fmt"
	"io"
	"path"
	"os"
	h "server/util/header"
)

//*Send File to client
func (s *Server) SendFileToClient(f string, client *Client) error {
	file, err := os.Open(f)
	
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Println(h.I, "Sending File to client:", client.ID)

	fileInfo, err := os.Stat(file.Name())
	if err != nil {
		return err
	}

	command := fmt.Sprintf("send %s %d", file.Name(), fileInfo.Size())
	
	
	if _, err := client.Conn.Write([]byte(command)); err != nil {
		fmt.Println(h.E, "Error sending command")
		return err
	}
	if _, err = io.CopyN(client.Conn, file, fileInfo.Size()); err != nil{
		fmt.Println(h.E, "Error copying file")
		return err
	}

	return nil
}

//*Sends command to prep download of file from client
func (s *Server) RecvFileFromClient(client *Client, f string) error {
	command := fmt.Sprintf("download %s", f)
	
	if _, err := client.Conn.Write([]byte(command)); err != nil{
		fmt.Println(h.E, "Error sending command")
		return err		
	}

	return nil
}

//*Downloads a file from the client
func (s *Server) DownloadFileFromClient(size int64, name string, c *Client) error{
	
	file, err := os.Create(path.Base(name))
	if err != nil{
		fmt.Println(h.E, "Error creating file locally")
		return err
	}
	defer file.Close()


	if _, err := io.CopyN(file, c.Conn, size); err != nil{
		fmt.Println(h.E, "Error copying file")
		return err
	}

	fmt.Println(h.K, "File downloaded")
	return nil
}