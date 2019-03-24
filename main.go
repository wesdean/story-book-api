package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"log"
	"net/http"
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

	r := mux.NewRouter()
	r.Handle("/", alice.New().ThenFunc(controllers.HealthCheckController{}.Index))

	httpPort := os.Getenv("HTTP_PORT")
	log.Printf("Server listening on port %s\n", httpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), r)
	if err != nil {
		log.Fatal("Error starting HTTP server")
	}
}
