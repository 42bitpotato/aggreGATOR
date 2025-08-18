package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

type State struct {
	cfg *config.Config
}

func main() {
	var cfg config.Config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}
	err = config.SetUser(&cfg, "Macke")
	if err != nil {
		fmt.Println("error setting user:", err)
		os.Exit(1)
	}
	cfg2, _ := config.Read()
	cfgJson, _ := json.Marshal(cfg2)
	fmt.Println(string(cfgJson))
}
