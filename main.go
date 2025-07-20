package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robgilliam/gator/internal/command"
	"github.com/robgilliam/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	state := command.State{
		Config: cfg,
	}

	commands := command.Commands{
		Commands: make(map[string]func(*command.State, command.Command) error),
	}

	commands.Register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: %s <command> [<args>]\n", args[0])
		os.Exit(1)
	}

	cmd := command.Command{
		Name: args[1],
		Args: args[2:],
	}

	if err = commands.Run(&state, cmd); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
