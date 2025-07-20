package command

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	return c.Commands[cmd.Name](s, cmd)
}
