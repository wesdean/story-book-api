package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartServer(httpPort string) *http.Server {
	r := mux.NewRouter()
	BindRoutes(r)

	log.Printf("Server listening on port %s\n", httpPort)
	server := &http.Server{Addr: fmt.Sprintf(":%s", httpPort), Handler: r}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting HTTP server")
	}
	return server
}
