package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/robgilliam/gator/internal/config"
	"github.com/robgilliam/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: gator <command> [args...]\n")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
