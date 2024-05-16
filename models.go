package main

import (
	"time"

	"github.com/Gambor27/RSSFeed/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIkey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		APIkey:    user.ApiKey,
	}
}

func databaseFeedToFeed(user database.Feed) Feed {
	return Feed{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Url:       user.Url,
		UserID:    user.UserID,
	}
}

func databaseFollowtoFollow(follow database.FeedsUser) FeedFollow {
	return FeedFollow{
		ID:        follow.ID,
		CreatedAt: follow.CreatedAt,
		UpdatedAt: follow.UpdatedAt,
		FeedID:    follow.FeedID,
		UserID:    follow.UserID,
	}
}
