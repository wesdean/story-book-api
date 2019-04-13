package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/wesdean/story-book-api/api"
	"log"
	"os"
)

func main() {
	var err error

	env := os.Getenv("ENV")
	if env == "" {
		env = "testing"
	}
	err = godotenv.Load(fmt.Sprintf("%s.env", env))
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	api.StartServer(os.Getenv("HTTP_PORT"))
}
