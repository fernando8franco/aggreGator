package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fernando8franco/aggreGator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	userName := cmd.Arguments[0]

	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("The user has been set")
	return nil
}

func HandlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	userName := cmd.Arguments[0]
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	}
	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created succesfully:")
	printUser(user)

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" - ID:   %v\n", user.ID)
	fmt.Printf(" - Name: %v\n", user.Name)
}
