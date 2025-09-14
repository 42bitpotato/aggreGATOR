package aggregator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/42bitpotato/aggreGATOR/internal/config"
	"github.com/42bitpotato/aggreGATOR/internal/database"
	"github.com/42bitpotato/aggreGATOR/internal/rss"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	for _, post := range feed.Channel.Item {

		var pubDate sql.NullTime
		parsedDate, err := dateparse.ParseAny(post.PubDate)
		if err != nil {
			s.Logger.Printf("could not parse pubDate %q: %v", post.PubDate, err)
			pubDate = sql.NullTime{Valid: false}
		} else {
			pubDate = sql.NullTime{Time: parsedDate, Valid: true}
		}

		postArgs := database.CreatePostParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Title:        post.Title,
			Url:          post.Link,
			Description:  post.Description,
			PublishedAt:  pubDate,
			PublishedRaw: post.PubDate,
			FeedID:       nextFeed.ID,
		}

		// Retry 3 times if error
		for i := 0; i < 3; i++ {
			err = s.Db.CreatePost(context.Background(), postArgs)
			if err != nil {
				retry, err := handleCreatePostErr(err)
				if retry {
					continue
				} else if err != nil {
					s.Logger.Printf("error creating post: %v", err)
				}
			}
			break
		}
	}

	//fmt.Printf("----------\nFETCHED FEED: %s\nITEMS:\n", feed.Channel.Title)
	//for _, item := range feed.Channel.Item {
	//	fmt.Printf("*	%s\n", item.Title)
	//}
	//fmt.Print("\n")

	return nil
}

func fetchFeed(url string) (feed *rss.RSSFeed, err error) {
	rssCli := rss.NewClient()

	Feed, err := rssCli.FetchFeed(context.Background(), url)
	if err != nil {
		return &rss.RSSFeed{}, err
	}

	return Feed, nil
}

func handleCreatePostErr(in error) (retry bool, out error) {
	if in == nil {
		return false, nil
	}
	var pqErr *pq.Error
	code := string(pqErr.Code)
	if errors.As(in, &pqErr) {
		if code == "23505" && pqErr.Constraint == "posts_url_key" {
			return false, nil
		} else if strings.HasPrefix(code, "08") {
			return true, in
		}
	}
	return false, in
}
