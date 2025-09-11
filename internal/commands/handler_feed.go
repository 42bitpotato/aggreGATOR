package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *config.State, cmd Command) error {
	args := cmd.Args
	if len(args) < 2 {
		return fmt.Errorf("2 arguments needed, name of feed and url: %v", args)
	}

	userId, err := getUserId(s)
	if err != nil {
		return err
	}

	feed := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    userId,
	}

	err = s.Db.AddFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	// dbFeed, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	// if err != nil {
	//	return fmt.Errorf("failed to fetch feed from db: %v", err)
	//}

	// newFollowArgs := database.CreateFeedFollowParams{
	//	ID:        uuid.New(),
	//	CreatedAt: time.Now(),
	//	UpdatedAt: time.Now(),
	//	UserID:    feed.UserID,
	//	FeedID:    feed.ID,
	//}

	followCmd := Command{
		Name: "follow",
		Args: cmd.Args[1:],
	}
	err = HandlerFollowFeed(s, followCmd)
	// folFeed, err := s.Db.CreateFeedFollow(context.Background(), newFollowArgs)
	if err != nil {
		return fmt.Errorf("failed to follow new feed: %v", err)
	}

	return nil
}

func HandlerResetFeeds(s *config.State, cmd Command) error {
	err := s.Db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset feeds: %v", err)
	}
	fmt.Println("feeds table reset.")
	return nil
}

func HandlerGetFeeds(s *config.State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve feeds from database: %v", err)
	}

	for _, item := range feeds {
		usrName, err := s.Db.GetUserByID(context.Background(), item.UserID)
		if err != nil {
			return fmt.Errorf("failed to retrieve username from database: %v", err)
		}
		printFeed := fmt.Sprintf("--- %s ---\nURL: %s\nCreated by: %s", item.Name, item.Url, usrName)

		if usrName == s.Cfg.CurrentUserName {
			printFeed += " (current)"
		}
		fmt.Printf("%s\nCreated: %s - Updated: %s\n\n",
			printFeed,
			item.CreatedAt.Format("2006-01-02 15:04:05"),
			item.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func HandlerFollowFeed(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("only 1 arguments needed, url: %v", cmd.Args)
	}
	url := cmd.Args[0]

	// Get feed row
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed through url: %v", err)
	}

	userID, err := getUserId(s)
	if err != nil {
		return fmt.Errorf("failed to fetch user ID: %v", err)
	}

	newFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	}

	newFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), newFollowArgs)
	if err != nil {
		return fmt.Errorf("error querrying sql database: %v", err)
	}

	fmt.Printf("User: %s\nFollowing: %s\n", newFeedFollow.UserName, newFeedFollow.FeedName)
	return nil
}

func HandlerUserFollowing(s *config.State, cmd Command) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetching user feeds: %v", err)
	}

	// Check number of feeds and adjust print accordingly
	numFeeds := len(feedFollows)
	printUsrF := fmt.Sprintf("User %s is following %v feed", s.Cfg.CurrentUserName, numFeeds)
	if numFeeds == 0 {
		fmt.Printf("%ss.\n", printUsrF)
		return nil
	}
	if numFeeds == 1 {
		fmt.Printf("%s:\n", printUsrF)
	} else if numFeeds > 1 {
		fmt.Printf("%ss:\n", printUsrF)
	}

	// Print followed feeds
	for _, feed := range feedFollows {
		fmt.Printf("- %s\n", feed.FeedName)
	}

	return nil
}
