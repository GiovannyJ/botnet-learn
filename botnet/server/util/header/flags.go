package header

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
)

const (
	I    = "[*]"
	E    = "[-]"
	K    = "[+]"
	L    = "[#]"
	X 	 = "[!]"
	Line = "[----------------------------------------------]"
)

//*List command menu
func ListHelp() {
	fmt.Println("[----------------------------------------------]")
	fmt.Println("[++++++++++++ STATUS COMMANDS ++++++++++++]")
	fmt.Printf("[*] count\t\t\t\t\t\t\t\t:: displays the number of active clients connected\n")
	fmt.Printf("[*] list\t\t\t\t\t\t\t\t:: displays active clients information\n")
	fmt.Printf("[*] metadata [-g active client group] [-a active client]\t\t:: sends command to return metadata of client\n")
	fmt.Println(I,"echo \t\t\t\t\t\t\t\t:: sends command that actively checks all clients to see if they're active")
	fmt.Println(I, "whoami [-g active client group] [-a active client]\t\t\t\t\t\t\t:: sends command to show what user is logged in on client side session")

	fmt.Println("[++++++++++++ CLIENT CONTROL ++++++++++++]")
	fmt.Printf("[*] set active <NUM>\t\t\t\t\t\t\t:: sets an active client to interact with (if set to 0 = global mode)\n")
	fmt.Printf("[*] check active\t\t\t\t\t\t\t:: displays the current active client ID\n")
	fmt.Printf("[*] set group <NUM>\t\t\t\t\t\t\t:: assigns random group of <NUM> size to be used in group mode\n")
	fmt.Printf("[*] check group\t\t\t\t\t\t\t\t:: displays all client IDs in active group\n")

	fmt.Println("[++++++++++++ CLIENT COMMANDS ++++++++++++]")
	fmt.Printf("[*] ping [-g active client group] [-a active client] <IPADDR>\t\t:: sends a command to ping <IPADDR> on client(s)\n")
	fmt.Printf("[*] entry [-g active client group] [-a active client]\t\t\t:: sends command to retrieve the PWD of the client\n")
	fmt.Printf("[*] search [-g active client group] [-a active client] <FILENAME>\t:: sends command to search for <FILENAME> on clients(s) (filepath must be abs)\n")
	fmt.Printf("[*] send [-g active client group] [-a active client] <FILENAME>\t\t:: sends command to send <FILENAME> from localhost to clients(s) (filepath must be abs)\n")
	fmt.Printf("[*] download [-g active client group] [-a active client] <FILENAME>\t:: sends command to download <FILENAME> from client to localhost (filepath must be abs)\n")
	fmt.Printf("[*] run [-g active client group] [-a active client] <FILENAME>\t\t:: sends command to run <FILENAME> on clients(s) (filepath must be abs)\n")
	fmt.Printf("[*] blowup [-g active client group] [-a active client]\t\t\t:: sends command to close the connection of client to server\n")
	fmt.Println("[++++++++++++ UTIL ++++++++++++]")
	fmt.Printf("[*] help/h\t\t\t\t:: displays this message\n")
	fmt.Printf("[*] exit\t\t\t\t:: exits the current session\n")
	fmt.Println(I,"clear\t\t\t\t:: clears the screen")
	fmt.Println("[----------------------------------------------]")
}


var clear map[string]func()

//*initialization of the clear map
func init() {
    clear = make(map[string]func())
    clear["linux"] = func() { 
        cmd := exec.Command("clear") 
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

//*method usd to clear the screen in it server session
func ClearScreen() {
    value, ok := clear[runtime.GOOS] 
    if ok { 
        value()  
    } else { 
        panic("Your platform is unsupported! I can't clear terminal screen :(")
    }
}

