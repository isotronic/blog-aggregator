package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/isotronic/blog-aggregator/internal/database"
	"github.com/lib/pq"
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
		publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
				// If parsing fails, use current time as fallback
				publishedAt = time.Now()
		}
		newPost := database.CreatePostParams{
			ID: uuid.New(),
			Title: sql.NullString{
				String: item.Title,
				Valid: true,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: publishedAt,
			FeedID: nextFeed.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.db.CreatePost(context.Background(), newPost)
		if err != nil {
			// Check if error is a unique violation on URL
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" && pqErr.Constraint == "posts_url_key" {
				// Ignore duplicate URL errors
				continue
			}
			fmt.Printf("Error creating post: %v\n", err)
		}
	}

	return nil
}
