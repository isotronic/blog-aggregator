package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/isotronic/blog-aggregator/internal/database"
)

func scrapeFeed(s *state) error {
	currentTime := sql.NullTime{
		Time: time.Now(), 
		Valid: true,
	}
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background(), currentTime)
	if err != nil {
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	markFetched := database.MarkFeedFetchedParams{
		LastFetchedAt: currentTime,
		ID: nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), markFetched)
	if err != nil {
		return err
	}

	fmt.Printf("Fetched feed: %v\n", nextFeed.Name)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("  * %v\n", item.Title)
	}

	return nil
}