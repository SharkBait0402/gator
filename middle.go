package main

import(
	"log"
	"context"
	"github.com/sharkbait0402/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {

		user, err:=s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err!=nil {
			log.Printf("couldn't get current user: %w\n", err)
			return err
		}

		return handler(s, cmd, user)

	}

}
