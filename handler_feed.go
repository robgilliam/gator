package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robgilliam/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed '%s': %w", cmd.Args[0], err)
	}

	fmt.Printf("created new feed: %+v", feed)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed follow for current user '%s' and url '%s': %w", user.Name, feed.Url, err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feedlist, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feed list: %w", err)
	}

	for _, feed := range feedlist {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("can't get user details for ID %s, feed %s: %w", feed.UserID, feed.Name, err)
		}

		fmt.Printf("* %s -> %s (%s)\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("no feed for url '%s': %w", cmd.Args[0], err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed follow for current user '%s' and url '%s': %w", user.Name, feed.Url, err)
	}

	fmt.Printf(" * Feed:    %s\n", feed.Name)
	fmt.Printf(" * Name:    %s\n", user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("can't get follows for user '%s': %w", user.Name, err)
	}

	for _, followed_feed := range following {
		fmt.Printf("%s\n", followed_feed.FeedName)
	}

	return nil
}
