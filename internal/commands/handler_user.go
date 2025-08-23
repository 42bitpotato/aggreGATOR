package commands

import (
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

func handlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing argument, the login handler expects a single argument, the username")
	}
	username := cmd.Args[0]
	err := config.SetUser(s.Cfg, username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to %s\n", username)
	return nil
}
