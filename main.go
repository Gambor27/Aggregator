package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	serverSetup(os.Getenv("PORT"))
}
