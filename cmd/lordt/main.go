package main

import (
	"github.com/jessehorne/lordt/internal/cli"
	"os"
)

func main() {
	ch := cli.NewCommandHandler()
	cli.InitCommands(ch)
	err := ch.Handle(os.Args)
	if err != nil {
		ch.RunCommand("help", []string{})
	}
}
