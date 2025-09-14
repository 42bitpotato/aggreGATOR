package commands

import (
	"fmt"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/aggregator"
	"github.com/42bitpotato/aggreGATOR/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommands map[string]func(*config.State, Command) error
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not supported: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.RegisteredCommands[name] = f
}

// Aggregator - Feed fetching automation
func Agg(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("agg takes 1 argument, time duration (1s, 1m, 1h): %v", cmd.Args)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time argument: %v", err)
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenReqs.Round(time.Second).String())

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		aggregator.ScrapeFeeds(s)
	}
}
