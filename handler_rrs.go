package main

import (
	"context"
	"fmt"
)

func HandlerAgg(s *state, cmd command) error {
	feed, err := fecthFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error in fetching the feed: %w", err)
	}

	printAggFeed(feed)

	return nil
}
