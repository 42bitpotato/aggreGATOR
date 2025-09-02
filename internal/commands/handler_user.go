package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing argument, the login handler expects a single argument, the username")
	}
	username := cmd.Args[0]

	// Check if user in database
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error checking username in database: %v", err)
	}

	// Set user
	err = config.SetUser(s.Cfg, username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to %s\n", username)
	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register command must contain exactly 1 argument: username. Input: %v", cmd.Args)
	}

	uName := cmd.Args[0]

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      uName,
	}
	usr, err := s.Db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to create new user in database: %v", err)
	}
	fmt.Printf("User '%s' successfully created.\n%v\n", uName, usr)

	// Login created user
	err = HandlerLogin(s, cmd)
	if err != nil {
		return fmt.Errorf("unable to set username in config: %v", err)
	}
	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users: %v", err)
	}
	fmt.Println("Successfully reseted users table!")
	return nil
}
