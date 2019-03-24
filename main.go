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

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
