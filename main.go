package main

import (
	"fmt"
	"os"

	"github.com/42bitpotato/aggreGATOR/internal/commands"
	"github.com/42bitpotato/aggreGATOR/internal/config"
)

// Get input
func getInput() (commands.Command, error) {
	if len(os.Args[:]) < 3 {
		return commands.Command{}, fmt.Errorf("atleast 2 arguments needed")
	}
	inputCmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	return inputCmd, nil
}

func main() {
	var state config.State
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}
	state.Cfg = &cfg

	cmds := commands.Commands{
		RegisteredCommands: make(map[string]func(*config.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)

	inputCmd, err := getInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmds.Run(&state, inputCmd)

	fmt.Print(cfg)
}
