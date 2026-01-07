package main

import (
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
	"os"
	"log"
	"github.com/sharkbait0402/gator/internal/database"
	"database/sql"
	"errors"
	"strconv"
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

	if len(cmd.args) < 1 {
		log.Println("no time bewtween req given")
		os.Exit(1)
	}

	duration, err:=time.ParseDuration(cmd.args[0])
	if err!=nil {
		log.Println("improper time format")
		os.Exit(1)
	}

	fmt.Printf("Collecting feeds every %v\n", cmd.args[0])

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	// url:="https://www.wagslane.dev/index.xml"
	//
	// feed, err:= fetchFeed(context.Background(), url)
	// if err!=nil {
	// 	log.Println("failed to get feed")
	// 	return err
	// }
	//
	// fmt.Println(*feed)
	//
	return nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 2 {
			log.Println("not enough args were given")
			os.Exit(1)
	}

	id:=uuid.New()
	now:=time.Now()
	name:=cmd.args[0]
	url:=cmd.args[1]
	
	feedParams:=database.CreateFeedParams {
		ID: id,
		CreatedAt: now,
		UpdatedAt: now,
		Name: name,
		Url: url,
		UserID: user.ID,
	}

	feed, err:=s.db.CreateFeed(context.Background(), feedParams)
	if err!=nil {
		log.Printf("failed to create feed: %w\n", err)
		return err
	}

	followParams:=database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	followRow, err:=s.db.CreateFeedFollow(context.Background(), followParams)
	if err!=nil {
		log.Println("failed to create followRow")
		return err
	}

	log.Println(feed)
	log.Printf("Followed %v as %v", followRow.FeedName, followRow.UserName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {

	feeds, err:=s.db.GetFeeds(context.Background())
	if err!= nil{
		log.Println(err)
		return err
	}

	for _, feed:=range feeds {
		fmt.Printf("* %v\n", feed.Name)
		fmt.Printf("* %v\n", feed.Url)
		fmt.Printf("* %v\n\n", feed.Name_2)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) == 0 {
		log.Println("no url was given")
		os.Exit(1)
	}

	url:=cmd.args[0]

	feed, err:=s.db.GetFeed(context.Background(), url)
	if err!=nil {
		log.Println("failed to get feed or feed does not exist")
		return err
	}

	params:=database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	followRow, err:=s.db.CreateFeedFollow(context.Background(), params)
	if err!=nil {
		log.Println("failed to create followRow")
		return err
	}

	log.Printf("Followed %v as %v", followRow.FeedName, followRow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	follows, err:=s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err!=nil {
		log.Printf("could not get users follow feed: %w\n", err)
		return err
	}

	for _, feed:=range follows {
		log.Printf("* %v\n",  feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		log.Println("no url was given")
		os.Exit(1)
	}

	feed,err:=s.db.GetFeed(context.Background(), cmd.args[0])
	if err!=nil {
		log.Println(err)
		return err
	}

	params:=database.UnfollowFeedForUserParams {
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err=s.db.UnfollowFeedForUser(context.Background(), params)
	if err!=nil {
		log.Println(err)
		return err
	}

	return err

}

func handlerBrowse(s *state, cmd command, user database.User) error {

	var limit int32
	limit = 2

	if len(cmd.args) >=1 {
		n, err:=strconv.Atoi(cmd.args[0])
		if err!=nil {
			log.Println("failed to convert arg to int: ", err)
		} else {
			limit = int32(n)
		}
	}

	params:=database.GetPostsParams {
		UserID: user.ID,
		Limit: limit,
	}

	posts, err:=s.db.GetPosts(context.Background(), params)
	if err!=nil {
		log.Println("falied to get post: ", err)
		return err
	}


	for _, post:=range posts {
		fmt.Printf("* %v\n", post.Title)
		fmt.Printf("* %v\n", post.Url)
		fmt.Printf("* %v\n", post.Description)
		fmt.Printf("* %v\n", post.PublishedAt)
		fmt.Println("\n")
	}

	return nil
}
