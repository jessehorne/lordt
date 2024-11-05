package cli

import (
	"fmt"
)

func HelpCommandHandler(ch *CommandHandler, args []string) error {
	fmt.Println(Logo, "\n")
	fmt.Println(Headline)
	fmt.Printf("Version: %s\n\n", Version)

	fmt.Println("Tools\n")
	for _, cmd := range ch.Commands {
		fmt.Printf("%s -\t%s\n", cmd.Name, cmd.Desc)
	}

	fmt.Println("\nTry `lordt <tool> help` to get more details on a specific tool.")
	return nil
}
