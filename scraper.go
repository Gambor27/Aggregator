package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Gambor27/RSSFeed/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping %v targets with a frequency of %v", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetFeedsToFetch(context.TODO(), int32(concurrency))
		if err != nil {
			log.Println("Failed to get feeds from DB:", err)
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
	err := db.MarkFetched(context.TODO(), feed.ID)
	if err != nil {
		log.Println("Error marking Feed fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error gathering feed from URL: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found Post: ", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found.", feed.Name, len(rssFeed.Channel.Item))
}
