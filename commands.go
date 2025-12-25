package main

import (
	"github.com/sharkbait0402/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handler map[string]func(*state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {

	handler, ok :=  c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {

	c.handlers[name] = f

}
