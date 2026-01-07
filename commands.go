package main

import (
	"github.com/sharkbait0402/gator/internal/config"
	"fmt"
	"github.com/sharkbait0402/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {

	handlers, ok :=  c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command %s", cmd.name)
	}

	return handlers(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {

	c.handlers[name] = f

}
