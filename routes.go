package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/middlewares"
)

func BindRoutes(r *mux.Router) []*mux.Route {
	return []*mux.Route{
		r.Handle("/", alice.New().ThenFunc(controllers.HealthCheckController{}.Index)),
		r.Handle("/authentication", alice.New(middlewares.DatabaseMiddleware).ThenFunc(controllers.AuthenticationController{}.CreateToken)).Methods("POST"),
	}
}
