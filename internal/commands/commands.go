package commands

import (
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	registeredCommands map[string]func(*config.State, Command) error
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not supported: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.registeredCommands[name] = f
}
