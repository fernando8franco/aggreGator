package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fernando8franco/aggreGator/internal/database"
	"github.com/google/uuid"
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

const UniqueUrlConstraint = "uq_posts_url"

func fecthFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	var feed RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("error decoding the body: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get the next feed: %w", err)
	}

	markFeedFetched := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
		ID:        nextFeed.ID,
	}
	_, err = s.db.MarkFeedFetched(context.Background(), markFeedFetched)
	if err != nil {
		return fmt.Errorf("couldn't get the mark the feed %s: %w", nextFeed.Name, err)
	}

	feedData, err := fecthFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("couldn't collect feed %s: %w", nextFeed.Name, err)
	}

	for _, item := range feedData.Channel.Item {
		pubDate, err := parseStringToTime(item.PubDate)
		if err != nil {
			return fmt.Errorf("couldn't parse the string %s: %w", item.PubDate, err)
		}

		newPost := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: sql.NullString{
				String: item.Title,
				Valid:  item.Title != "",
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Title,
				Valid:  item.Title != "",
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: item.PubDate != "",
			},
			FeedID: nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), newPost)
		if err != nil {
			if strings.Contains(err.Error(), UniqueUrlConstraint) {
				continue
			}
			log.Printf("couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(feedData.Channel.Item))

	return nil
}

func parseStringToTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}

	formats := []string{
		time.DateTime,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	}

	var t time.Time
	var err error
	for _, layout := range formats {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("no matching format found for %q", timeStr)
}
