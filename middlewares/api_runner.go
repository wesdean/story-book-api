package middlewares

import (
	"github.com/gorilla/context"
	"net/http"
)

func AppendCleanups(r *http.Request, cleanup func(h http.Handler) http.Handler) {
	var cleanups []func(h http.Handler) http.Handler
	cleanupsContext, ok := context.GetOk(r, "Cleanup")
	if ok {
		cleanups = cleanupsContext.([]func(h http.Handler) http.Handler)
	}
	cleanups = append(cleanups, cleanup)
	context.Set(r, "Cleanup", cleanups)
}

func RunAPI(action func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.HandlerFunc(action).ServeHTTP(w, r)
		CleanupMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		).ServeHTTP(w, r)
	})
}
