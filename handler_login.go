package main

import "fmt"

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username is required")
	}

	userName := cmd.Arguments[0]
	err := s.Conf.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil
}
