package api

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/middlewares"
)

func BindRoutes(r *mux.Router) []*mux.Route {
	return []*mux.Route{
		r.Handle("/", alice.New(middlewares.ConfigMiddleware, middlewares.LoggingMiddleware, middlewares.DatabaseMiddleware).Then(middlewares.RunAPI(controllers.HealthCheckController{}.Index))),
		r.Handle("/authentication", alice.New(middlewares.ConfigMiddleware, middlewares.LoggingMiddleware, middlewares.DatabaseMiddleware).Then(middlewares.RunAPI(controllers.AuthenticationController{}.CreateToken))).Methods("POST"),
		r.Handle("/user_roles", alice.New(middlewares.ConfigMiddleware, middlewares.LoggingMiddleware, middlewares.DatabaseMiddleware).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))).Methods("GET"),
	}
}
