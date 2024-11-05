package cli

import "errors"

var (
	ErrNoArgs         = errors.New("no args provided")
	ErrInvalidCommand = errors.New("invalid command")
)

type CommandHandler struct {
	Commands map[string]*Command
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		Commands: map[string]*Command{},
	}
}

func (ch *CommandHandler) AddCommand(c *Command) {
	ch.Commands[c.Name] = c
}

func (ch *CommandHandler) RunCommand(short string, args []string) error {
	cmd, ok := ch.Commands[short]
	if !ok {
		return ErrInvalidCommand
	}
	err := cmd.Handler(ch, args)
	if err != nil {
		return err
	}
	return nil
}

func (ch *CommandHandler) Handle(args []string) error {
	if len(args) < 2 {
		return ErrNoArgs
	}
	return ch.RunCommand(args[1], args[2:])
}

func InitCommands(ch *CommandHandler) {
	ch.AddCommand(&Command{
		Name:    "help",
		Desc:    "Display helpful information about the Lordt CLI.",
		Handler: HelpCommandHandler,
	})

	ch.AddCommand(&Command{
		Name:    "nlock",
		Desc:    "Lock a file or files from being read as long as this is running. (experimental, slow, unreliable)",
		Handler: NLockCommandHandler,
	})
}
