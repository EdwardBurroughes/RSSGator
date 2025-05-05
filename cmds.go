package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/EdBurroughes/rss-blog-aggregator/internal/config"
	"github.com/EdBurroughes/rss-blog-aggregator/internal/database"
	"github.com/EdBurroughes/rss-blog-aggregator/internal/feeds"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.cmds[name]; ok {
		fmt.Printf("command %s already registered\n", name)
		return
	}
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}

	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no arguments expecting a username")
	}
	name := cmd.arguments[0]
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user %s does not exist", name)
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %v", user)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("issue setting the username: %v", err)
	}
	fmt.Println("sucessfully set the user")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("incorrect argument: %v", cmd.arguments)
	}
	ctx := context.Background()
	name := cmd.arguments[0]
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("failed to set username to cfg: %v", err)
	}
	fmt.Printf("successfully created user: %v", user)
	return nil
}

func handleReset(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("expecting no arguments: %v", cmd.arguments)
	}
	ctx := context.Background()
	if err := s.db.DeleteUser(ctx); err != nil {
		return err
	}
	return nil
}

func handleGetUsers(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("expecting no arguments: %v", cmd.arguments)
	}
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		userStr := fmt.Sprintf("* %s", user)
		if user == *s.cfg.CurrentUserName {
			userStr += " (current)"
		}
		fmt.Printf("%s\n", userStr)
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	ctx := context.Background()

	duration, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", duration)
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := feeds.ScrapeFeeds(ctx, s.db)
		if err != nil {
			fmt.Printf("Found an error scraping: %v\n", err)
		}
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, *s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if length_args := len(cmd.arguments); length_args != 2 {
		return fmt.Errorf("uh oh expecting 2 arguments, %d provided: %v", length_args, cmd.arguments)
	}
	ctx := context.Background()

	// should be run in a transaction
	feed, err := s.db.CreateFeed(
		ctx,
		database.CreateFeedParams{
			Name:   cmd.arguments[0],
			Url:    cmd.arguments[1],
			UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		},
	)
	s.db.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed %v", err)
	}
	fmt.Printf("successfully created feed: %v", feed)
	return nil
}

func handleFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("*%s '%s' %s\n", feed.UserName, feed.FeedName, feed.Url)
	}
	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("should be 1 argument url, the following provided: %v", cmd.arguments)
	}
	ctx := context.Background()

	feed, err := s.db.GetFeed(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("user %s following %s", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Println("Following: ")
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed)
	}
	return nil
}

func handleUnfollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feed, err := s.db.GetFeed(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}
	unfollowed, err := s.db.DeleteFeedFollowForUser(ctx, database.DeleteFeedFollowForUserParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Unfollowed %s, id: %d", cmd.arguments[0], unfollowed.ID)
	return nil
}

func handleBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("should be 1 argument limit, the following provided: %v", cmd.arguments)
	}
	ctx := context.Background()

	limit, err := strconv.Atoi(cmd.arguments[0])
	if err != nil {
		return err
	}

	posts, err := s.db.GetPosts(ctx, database.GetPostsParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Println(post.Title.String)
	}
	return nil
}
