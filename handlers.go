package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isotronic/blog-aggregator/internal/database"
)

func aggHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing duration argument")
	}

	timeBetween, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetween.String())
	ticker := time.NewTicker(timeBetween)
	for ; ; <-ticker.C {
		err := scrapeFeed(s)
		if err != nil {
			fmt.Printf("Error fetching feed: %v\n", err) 
		}
	}
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing username")
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		Name: cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	newUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User created: %v\n", newUser)
	return nil
}

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing username")
	}

	user, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("user not found")
		}
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", user.Name)
	return nil
}

func usersHandler(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	loggedIn := ""
	fmt.Println("Users:")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			loggedIn = "(current)"
		}
		fmt.Printf(" * %v %v\n", user.Name, loggedIn)
		loggedIn = ""
	}

	return nil
}

func resetHandler(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Users successfully reset")
	return nil
}

func addFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing feed name or url")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	feed := database.CreateFeedParams{
		ID: uuid.New(),
		Name: name,
		Url: url,
		UserID: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	newFollow := database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: newFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		return err
	}

	fmt.Printf("Feed added: %v\n", newFeed)
	return nil
}

func feedHandler(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		userName, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			userName.Name = "Error fetching user"
		}
		fmt.Printf("  Name: %v\n", feed.Name)
		fmt.Printf("  URL: %v\n", feed.Url)
		fmt.Printf("  User: %v\n", userName.Name)
		fmt.Println("-----")
	}

	return nil
}

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing feed url")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	newFollow := database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(),newFollow)
	if err != nil {
		return err
	}

	fmt.Printf("%v followed feed: %v\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(),user.ID)
	if err != nil {
		return err
	}

	fmt.Println("You are following:")
	for _, feed := range feeds {
		fmt.Printf(" * %v\n", feed.FeedName)
	}

	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing feed url")
	}

	unfollow := database.DeleteFeedFollowByUrlParams{
		UserID: user.ID,
		Url: cmd.args[0],
	}
	err := s.db.DeleteFeedFollowByUrl(context.Background(), unfollow)
	if err != nil {
		return err
	}
	
	fmt.Printf("You unfollowed feed: %v\n", cmd.args[0])

	return nil
}

func browseHandler(s *state, cmd command, user database.User) error {
	limit := int32(2)
	offset := int32(0)
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	}
	if len(cmd.args) > 1 {
		parsedOffset, err := strconv.ParseInt(cmd.args[1], 10, 32)
		if err != nil {
			return err
		}
		offset = int32(parsedOffset)
	}
	
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit,
		Offset: offset,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Posts:\n")
	for _, post := range posts {
		fmt.Printf(" * %v\n", post.Title.String)
	}

	return nil
}
