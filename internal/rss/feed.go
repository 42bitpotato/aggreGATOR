package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/42bitpotato/aggreGATOR/internal/config"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *Client) FetchFeed(s *config.State, ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Set header
	req.Header.Set("User-Agent", "gator")

	// Get respons
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	// Read response (return byte)
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Decode/Unmarshall response
	rssResp := &RSSFeed{}
	err = xml.Unmarshal(dat, &rssResp)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Unescape respons
	err = unescapeHTML(s, rssResp)
	if err != nil {
		return &RSSFeed{}, err
	}

	return rssResp, nil
}

func unescapeHTML(s *config.State, feed *RSSFeed) error {
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	// Strip HTML code
	feed.Channel.Description = s.HTMLpolicy.Sanitize(feed.Channel.Description)

	if len(feed.Channel.Item) == 0 {
		return nil
	}
	for i, item := range feed.Channel.Item {
		itemDesc := html.UnescapeString(item.Description)
		itemTitel := html.UnescapeString(item.Title)

		// Handel description
		stripDescription := s.HTMLpolicy.Sanitize(itemDesc)
		splitDesc := strings.Split(stripDescription, "\n")
		newDesc := ""
		for j, line := range splitDesc {
			lineFix := strings.TrimSpace(line)
			lineFix = strings.ToLower(lineFix)
			if lineFix == "" {
				continue
			}
			if strings.HasPrefix(lineFix, "article url") || strings.HasPrefix(lineFix, "comments url") {
				continue
			}
			if strings.HasPrefix(lineFix, "points:") || strings.HasPrefix(lineFix, "# comments") {
				continue
			}
			newDesc += line
			if j < len(splitDesc)-1 {
				newDesc += fmt.Sprint("\n")
			}
		}

		// If description is empty, trim before return empty string
		if strings.TrimSpace(newDesc) == "" {
			feed.Channel.Item[i].Description = ""
		} else {
			feed.Channel.Item[i].Description = newDesc

		}
		feed.Channel.Item[i].Title = itemTitel
	}
	return nil
}
