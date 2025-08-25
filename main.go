package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/fernando8franco/aggreGator/internal/config"
	"github.com/fernando8franco/aggreGator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", conf.DBUrl)
	if err != nil {
		log.Fatalf("error in the database connection: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := state{
		cfg: &conf,
		db:  dbQueries,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	commands.Register("login", HandlerLogin)
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	commands.Register("users", HandlerUsers)
	commands.Register("agg", HandlerAgg)
	commands.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	commands.Register("feeds", HandlerFeeds)
	commands.Register("follow", middlewareLoggedIn(HandlerFollow))
	commands.Register("following", middlewareLoggedIn(HandlerFollowing))
	commands.Register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	commands.Register("browse", middlewareLoggedIn(HandleBrowse))

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments were provided")
	}

	command := command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}

	err = commands.Run(&programState, command)
	if err != nil {
		log.Fatal(err)
	}
}
