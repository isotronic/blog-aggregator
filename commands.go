package main

type command struct {
	name  string
	args  []string
}

type commands map[string]func(*state, command) error

func (c *commands) register(name string, f func(*state, command) error) {
	(*c)[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	return (*c)[cmd.name](s, cmd)
}