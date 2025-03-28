package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isotronic/blog-aggregator/internal/database"
)

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

func aggHandler(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func addFeedHandler(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing feed name or url")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

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