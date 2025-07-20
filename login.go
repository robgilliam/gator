package main

import (
	"fmt"

	"github.com/robgilliam/gator/internal/command"
)

func handlerLogin(s *command.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("%s requires 1 argument", cmd.Name)
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.Args[0])

	return nil
}
