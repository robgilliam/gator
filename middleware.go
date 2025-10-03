package main

import (
	"context"
	"fmt"

	"github.com/robgilliam/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("couldn't get details for current user %s: %w", s.cfg.CurrentUsername, err)
		}

		return handler(s, cmd, user)
	}
}
