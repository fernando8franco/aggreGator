package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fernando8franco/aggreGator/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %v \"<name>\" \"<url>\"", cmd.Name)
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
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created succesfully:")
	printFeed(feed)
	fmt.Println("Feed Follow created succesfully:")
	fmt.Printf("* User Name:     %s\n", feedFollow.UserName)
	fmt.Printf("* Feed Name:     %s\n", feedFollow.FeedName)

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

func HandleBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("usage: %v <limit>", cmd.Name)
	}

	limit := 2
	if len(cmd.Arguments) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Arguments[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	getPostsForUserParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostsForUserParams)
	if err != nil {
		return fmt.Errorf("couldn't get the posts for user: %w", err)
	}

	for _, post := range posts {
		printPost(post)
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

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("* ID:                %s\n", post.ID)
	fmt.Printf("* Title:             %s\n", post.Title.String)
	fmt.Printf("* Url:               %s\n", post.Url)
	fmt.Printf("* Published At:      %s\n", post.PublishedAt.Time)
	fmt.Printf("* Feed Name:         %s\n", post.FeedName)
	fmt.Println("================================================")
}
