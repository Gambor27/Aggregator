package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	serverSetup(os.Getenv("PORT"))
}
