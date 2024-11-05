package cli

type Command struct {
	Name    string
	Desc    string
	Handler func(*CommandHandler, []string) error
}
