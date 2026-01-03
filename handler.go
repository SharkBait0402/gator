package main

import (
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
	"os"
	"log"
	"github.com/sharkbait0402/blog-aggregator/internal/database"
	"database/sql"
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	
	if len(cmd.args) == 0 {
		log.Println("no username was given")
		return fmt.Errorf("no username was given")
	}

	_, err:=s.db.GetUser(context.Background(), cmd.args[0])
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("User does not exist")
		return err
	} else if err!=nil {
		fmt.Println("failed to retrieve user")
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("username was set to %v\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	
	if len(cmd.args) == 0 {
		log.Println("no name was given")
		return fmt.Errorf("no name")
	}

	ctx:=context.Background()
	id:=uuid.New()
	now:=time.Now()
	name:=cmd.args[0]

	params:= database.CreateUserParams{
		ID: id,
		CreatedAt: now,
		UpdatedAt: now,
		Name: name,
	}


	user, err:=s.db.GetUser(ctx, name)
	if err!=nil {
		if errors.Is(err, sql.ErrNoRows) {

			user, err= s.db.CreateUser(ctx, params)
				if err!=nil {
					return fmt.Errorf("create user failed: %w", err)
				}

		} else {
			return fmt.Errorf("get user failed: %w", err)
		}
	} else {
		os.Exit(1)
	}

	log.Printf("created user: %+v\n", user)

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("username was set to %v\n", cmd.args[0])


	return nil

}

func handlerReset(s *state, cmd command) error {

	err:=s.db.Reset(context.Background())
	if err!=nil {
		log.Println("error resetting database")
		return err
	}

	fmt.Println("database reset successfully")

	return nil
}

func handlerUsers(s *state, cmd command) error {

	names, err:=s.db.GetUsers(context.Background())
	if err!=nil {
		return err
	}

	for _, name:=range names {
		if s.cfg.CurrentUserName == name {
			fmt.Printf("* %v (current)\n", name)
			continue
		}
		fmt.Printf("* %v\n", name)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	url:="https://www.wagslane.dev/index.xml"
	
	feed, err:= fetchFeed(context.Background(), url)
	if err!=nil {
		log.Println("failed to get feed")
		return err
	}

	fmt.Println(*feed)

	return nil

}
