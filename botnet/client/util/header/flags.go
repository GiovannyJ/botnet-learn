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
	fmt.Println(Line)
}