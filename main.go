package main

import (
	"log"
	"os"

	"github.com/fernando8franco/aggreGator/internal/config"
)

type state struct {
	Conf *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := state{
		Conf: &conf,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	commands.Register("login", HandlerLogin)

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
