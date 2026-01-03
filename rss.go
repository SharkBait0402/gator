package main

import(
	"net/http"
	"encoding/xml"
	"html"
	"io"
	"context"
	"time"
	"log"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	client:=http.Client {
		Timeout: 10 * time.Second,
	}

	req, err:=http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err!=nil {
		log.Println("error with request")
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err:=client.Do(req)
	if err!=nil {
		return nil, err
	}
	defer resp.Body.Close()

	info, err:= io.ReadAll(resp.Body)
	if err!=nil {
		return nil, err
	}

	data:=&RSSFeed{}

	err=xml.Unmarshal(info, data)
	if err!=nil {
		return nil, err
	}

	data.Channel.Title=html.UnescapeString(data.Channel.Title)
	data.Channel.Description=html.UnescapeString(data.Channel.Description)

	for idx, item:=range data.Channel.Item {
		data.Channel.Item[idx].Title=html.UnescapeString(item.Title)
		data.Channel.Item[idx].Description=html.UnescapeString(item.Description)
	}

	return data, nil

}
