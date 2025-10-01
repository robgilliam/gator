package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robgilliam/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user '%s': %w", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user '%s': %w", name, err)
	}

	fmt.Printf("User switched successfully!\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("couldn't create user '%s': %w", name, err)
	}

	err = s.cfg.SetUser(name)

	if err != nil {
		return fmt.Errorf("couldn't set current user '%s': %w", name, err)
	}

	fmt.Printf("User created successfully:\n")
	printUser(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	if err := s.db.DropUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't reset user database: %w", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	userlist, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get user list: %w", err)
	}

	for _, user := range userlist {
		name := user.Name
		if name == s.cfg.CurrentUsername {
			name += " (current)"
		}

		fmt.Printf("* %s\n", name)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
