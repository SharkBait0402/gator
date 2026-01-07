package main

import(
	"context"
	"log"
	"fmt"
	"github.com/sharkbait0402/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"github.com/lib/pq"
	"database/sql"
)

func scrapeFeeds(s *state) {

	nextFeed,err:=s.db.GetNextFeedToFetch(context.Background())
	if err!=nil{
		log.Println("failed to get next feed to fetch")
	}

	err=s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err!=nil {
		log.Println("failed to mark feed")
	}

	feed, err:=fetchFeed(context.Background(), nextFeed.Url)
	if err!=nil {
		log.Println("failed to fetch feed")
	}

	for _, item:=range feed.Channel.Item {
		if item.Title == "" {
			continue
		}

		layout:=time.RFC1123Z

		pubDate,err:=time.Parse(layout, item.PubDate)
		if err!=nil {
			log.Println("error parsing pub date: ", err)
			continue
		}

		params:=database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: sql.NullTime{
				Time: pubDate,
				Valid: true,
			},
			FeedID: nextFeed.ID,
		}
		
		_,err=s.db.CreatePost(context.Background(), params)
		if err!=nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505"{
				continue
			} else {
				log.Println("error creating post: ", err)
				continue
			}

		}

		fmt.Printf("* %v\n", item.Title)
	}

	fmt.Println()

}
