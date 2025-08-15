package main

import (
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

func main() error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	err = config.SetUser(cfg, "Macke")
	if err != nil {
		return nil
	}
	cfg, err = config.Read()
	if err != nil {
		return err
	}
	fmt.Print(cfg)
	return nil
}
