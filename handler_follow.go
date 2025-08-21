package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fernando8franco/aggreGator/internal/database"
	"github.com/google/uuid"
)

func HandlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}

	url := cmd.Arguments[0]
	feedID, err := s.db.GetFeedID(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get the feed id: %w", err)
	}

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed Follow created succesfully:")
	fmt.Printf("* User Name:     %s\n", feedFollow.UserName)
	fmt.Printf("* Feed Name:     %s\n", feedFollow.FeedName)

	return nil
}

func HandlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}

	url := cmd.Arguments[0]
	deleteFeedFollow := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}
	err := s.db.DeleteFeedFollow(context.Background(), deleteFeedFollow)
	if err != nil {
		return fmt.Errorf("couldn't delete the feed follow from the database")
	}

	return nil
}

func HandlerFollowing(s *state, cmd command, user database.User) error {
	feedsFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get the feeds follows: %w", err)
	}

	for _, feedFollow := range feedsFollows {
		fmt.Printf("* Feed Name:     %s\n", feedFollow.FeedName)
	}

	return nil
}
