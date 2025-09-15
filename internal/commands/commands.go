package commands

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/aggregator"
	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
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

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs.Round(time.Second).String())

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		aggregator.ScrapeFeeds(s)
	}
}

func HandlerBrowse(s *config.State, cmd Command, user database.User) error {
	// Set return limit
	rLimit := 2
	// Handle argument errors
	if len(cmd.Args) == 0 {
		fmt.Println("Browse: No arguments given, return limit set to 2.")
	} else if len(cmd.Args) == 1 {
		num, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("argument needs to be only numbers: %v", cmd.Args[0])
		} else if num >= 1 && num <= 50 {
			fmt.Printf("Browse: Limit set to %v\n", num)
			rLimit = num
		} else {
			return fmt.Errorf("limit needs to be between 1 and 50: %v", num)
		}
	} else {
		return fmt.Errorf("browse command may contain max 1 argument, an optional 'limit' parameter only numbers: %v", cmd.Args)
	}

	args := database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(rLimit),
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("error fetching user posts from db: %v", err)
	}

	fmt.Printf("USER: %s\n", s.Cfg.CurrentUserName)
	fmt.Println("---LATEST POSTS---")
	for _, post := range posts {
		var pubDate string
		if post.PublishedAt.Valid {
			pubDate = post.PublishedAt.Time.Format(s.Cfg.DateFormat)
		} else {
			pubDate = ""
		}
		stripDescription := s.HTMLpolicy.Sanitize(post.Description)

		fmt.Printf("TITLE: %s\n", post.Title)
		fmt.Printf("*	Description: %s\n", stripDescription)
		fmt.Printf("*	URL: %s\n", post.Url)
		fmt.Printf("*	Published: %s\n\n", pubDate)
	}

	return nil
}
