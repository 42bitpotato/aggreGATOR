package commands

import (
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument, the login handler expects a single argument, the username")
	}
	username := cmd.args[0]
	err := config.SetUser(s.cfg, username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to %s\n", username)
	return nil
}
