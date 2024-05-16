package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Gambor27/RSSFeed/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var newUser = database.CreateUserParams{}
	uuid, err := uuid.NewUUID()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUser.ID = uuid
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	newUser.Name = request.Name

	user, err := cfg.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, databaseUserToUser(user))
}

func (cfg apiConfig) getUserByKey(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg apiConfig) followFeed(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request struct {
		Feed string `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedUUID, err := uuid.Parse(request.Feed)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newFollow = database.CreateFeedFollowParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedUUID,
		UserID:    user.ID,
	}

	follow, err := cfg.DB.CreateFeedFollow(r.Context(), newFollow)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, databaseFollowtoFollow(follow))
}

func (cfg apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	followID := r.PathValue("followID")
	followUUID, err := uuid.Parse(followID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	follow, err := cfg.DB.GetFollow(r.Context(), followUUID)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}

	if follow.UserID != user.ID {
		respondWithJSON(w, http.StatusUnauthorized, "Access denied")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), followUUID)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusAccepted, "Unfollowed")
}

func (cfg apiConfig) getUserFollows(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	follows, err := cfg.DB.GetUserFeeds(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusAccepted, follows)
}

func (cfg apiConfig) getPostsforUser(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getKey(r)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var params = database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  10,
	}
	posts, err := cfg.DB.GetPostByUser(r.Context(), params)
	if err != nil {
		jsonError(w, http.StatusNotFound, err.Error())
		return
	}
	var output = make([]PostByUser, 0)

	for _, post := range posts {
		nextPost := databasePostsByUserToPostByUser(post)
		output = append(output, nextPost)
	}
	respondWithJSON(w, http.StatusOK, output)
}
