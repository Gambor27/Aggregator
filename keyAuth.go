package main

import (
	"net/http"

	"github.com/Gambor27/RSSFeed/internal/database"
)

func (cfg apiConfig) getKey(r *http.Request) (database.User, error) {
	authHeader := r.Header.Get("Authorization")
	key := authHeader[7:]
	user, err := cfg.DB.GetUserByKey(r.Context(), key)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}
