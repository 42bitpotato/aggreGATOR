package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func HandlerFollowFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("only 1 arguments needed, url: %v", cmd.Args)
	}
	url := cmd.Args[0]

	// Get feed row
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed through url: %v", err)
	}

	newFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	newFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), newFollowArgs)
	if err != nil {
		return fmt.Errorf("error querrying sql database: %v", err)
	}

	fmt.Printf("User: %s\nFollowing: %s\n", newFeedFollow.UserName, newFeedFollow.FeedName)
	return nil
}

func HandlerUserFollowing(s *config.State, cmd Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error fetching user feeds: %v", err)
	}

	// Check number of feeds and adjust print accordingly
	numFeeds := len(feedFollows)
	printUsrF := fmt.Sprintf("User %s is following %v feed", user.Name, numFeeds)
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
