package aggregator

import (
	"context"
	"fmt"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/rss"
)

func ScrapeFeeds(s *config.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching next feed: %v", err)
	}

	err = s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("error marking feed with timestamp: %v", err)
	}

	feed, err := fetchFeed(nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching RSS feed: %v", err)
	}

	fmt.Printf("Fetched feed: %s\nItems:\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf("*	%s\n", item.Title)
	}
	fmt.Print("\n")

	return nil
}

func fetchFeed(url string) (feed *rss.RSSFeed, err error) {
	rssCli := rss.NewClient()

	var Feed *rss.RSSFeed
	Feed, err = rssCli.FetchFeed(context.Background(), url)
	if err != nil {
		return &rss.RSSFeed{}, err
	}

	return Feed, nil
}
