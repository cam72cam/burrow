package main

import (
	"fmt"
	"os"

	"github.com/cam72cam/burrow/commands"
	"github.com/cam72cam/burrow/display"
)

func main() {
	fn, err := display.Init()
	if err != nil {
		fmt.Println("Error initializing display: %v", err)
		os.Exit(1)
	}
	defer fn()

	for {
		in := display.NextInput()
		if in.String() == ":" {
			cmd := display.NextCommand(commands.MatchInput)
			if cmd != nil {
				commands.Run(nil, cmd.Name, cmd.Params)
			}
		} else {
			//TODO
		}
	}
}
