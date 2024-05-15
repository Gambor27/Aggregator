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
	var request struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authHeader := r.Header.Get("Authorization")
	key := authHeader[7:]
	user, err := cfg.DB.GetUserByKey(r.Context(), key)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
