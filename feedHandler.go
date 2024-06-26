package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Gambor27/RSSFeed/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) createFeed(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var newFeed = database.CreateFeedParams{}
	feedUUID, err := uuid.NewUUID()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newFeed.ID = feedUUID
	newFeed.CreatedAt = time.Now()
	newFeed.UpdatedAt = time.Now()
	newFeed.Name = request.Name
	newFeed.Url = request.URL
	newFeed.UserID = user.ID

	dbFeed, err := cfg.DB.CreateFeed(r.Context(), newFeed)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	followuuid, err := uuid.NewUUID()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feed := databaseFeedToFeed(dbFeed)

	var newFollow = database.CreateFeedFollowParams{
		ID:        followuuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	dbFollow, err := cfg.DB.CreateFeedFollow(r.Context(), newFollow)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	follow := databaseFollowtoFollow(dbFollow)

	var ResponseStruct = struct {
		Feed   Feed
		Follow FeedFollow
	}{
		Feed:   feed,
		Follow: follow,
	}
	respondWithJSON(w, http.StatusAccepted, ResponseStruct)
}

func (cfg apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	feeds := make([]Feed, 0)
	for _, feed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(feed))
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
