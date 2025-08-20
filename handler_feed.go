package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fernando8franco/aggreGator/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *state, cmd command) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %v \"<name>\" \"<url>\"", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get the user: %w", err)
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	fmt.Println("Feed created succesfully:")
	printFeed(feed)

	return nil
}

func HandlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get the feeds: %w", err)
	}

	fmt.Println("==============================================================")
	for _, feed := range feeds {
		fmt.Printf("* ID:            %s\n", feed.ID)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* User Name:     %s\n", feed.UserName)
		fmt.Println("==============================================================")
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
