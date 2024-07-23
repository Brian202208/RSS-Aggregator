package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Brian202208/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Print("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Panicln("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)

	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Couldn't parse date %v with err %v", item.PubDate, err)
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: publishedAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post:", err)
		}

	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
