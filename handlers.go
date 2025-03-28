package main

import (
	"context"
	"fmt"
	"log"
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

	log.Printf("User created: %v\n", newUser)
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

	log.Printf("Logged in as %s\n", user.Name)
	return nil
}