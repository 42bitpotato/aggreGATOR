package main

import (
	"fmt"
	"os"

	"github.com/42bitpotato/aggreGATOR/internal/commands"
	"github.com/42bitpotato/aggreGATOR/internal/config"
)

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

}
