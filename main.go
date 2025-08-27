package main

import (
	"fmt"
	"os"

	"github.com/42bitpotato/aggreGATOR/internal/commands"
	"github.com/42bitpotato/aggreGATOR/internal/config"
)

func encodeJson(conf config.Config) (string, error) {
	return "", nil
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

	inputCmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	cmds.Run(&state, inputCmd)

	fmt.Print(cfg)
}
