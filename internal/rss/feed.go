package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
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

func (c *Client) FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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
	err = xml.Unmarshal(dat, rssResp)
	if err != nil {
		return &RSSFeed{}, err
	}

	return &RSSFeed{}, nil
}
