package main

import(
	"context"
	"log"
	"fmt"
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
		fmt.Printf("* %v\n", item.Title)
	}

	fmt.Println()

}
