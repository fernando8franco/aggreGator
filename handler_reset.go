package main

import (
	"context"
	"fmt"
)

func HandlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete the users from the database")
	}

	return nil
}
