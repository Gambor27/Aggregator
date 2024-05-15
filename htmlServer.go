package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Gambor27/RSSFeed/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func serverSetup(port string) error {
	dbURL := os.Getenv(("CONNECTION_STRING"))
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err.Error())
	}
	var cfg = apiConfig{}
	cfg.DB = database.New(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/readiness", health)
	mux.HandleFunc("GET /v1/err", giveError)
	mux.HandleFunc("POST /v1/users", cfg.createUser)
	mux.HandleFunc("GET /v1/users", cfg.getUserByKey)
	mux.HandleFunc("POST /v1/feeds", cfg.createFeed)
	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)

	corsMux := middlewareCors(mux)
	address := fmt.Sprintf("localhost:%s", port)
	server := &http.Server{
		Addr:    address,
		Handler: corsMux,
	}
	err = server.ListenAndServe()
	if err != nil {
		return errors.New("server failed to start")
	}
	return nil
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func giveError(w http.ResponseWriter, r *http.Request) {
	jsonError(w, 500, "Internal Server Error")
}
