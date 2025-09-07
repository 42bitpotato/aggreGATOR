package commands

import (
	"context"
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/rss"
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

// Aggregator command, will be automated later on
func Agg(s *config.State, cmd Command) error {
	link := "https://www.wagslane.dev/index.xml"
	rssCli := rss.NewClient()
	feed, err := rssCli.FetchFeed(context.Background(), link)
	if err != nil {
		return fmt.Errorf("error fetching RSS feed: %v", err)
	}
	fmt.Print(feed)
	return nil
}
