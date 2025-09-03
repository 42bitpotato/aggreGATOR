package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/42bitpotato/aggreGATOR/internal/commands"
	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
	_ "github.com/lib/pq"
)

// Get input
func getInput() (commands.Command, error) {
	if len(os.Args[:]) < 2 {
		return commands.Command{}, fmt.Errorf("atleast 1 command needed")
	}
	inputCmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	return inputCmd, nil
}

func main() {
	// Load the config file
	var state config.State
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}
	state.Cfg = &cfg

	// Open a connection to the database, and store it in the state struct
	db, err := sql.Open("postgres", state.Cfg.DbUrl)
	if err != nil {
		fmt.Printf("error connectiong to sql database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state.Db = dbQueries

	// Generate commands
	cmds := commands.Commands{
		RegisteredCommands: make(map[string]func(*config.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerGetUsers)

	// Get input
	inputCmd, err := getInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run command
	err = cmds.Run(&state, inputCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
