package middlewares

import (
	"github.com/gorilla/context"
	"net/http"
)

func CleanupMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cleanups []func(h http.Handler) http.Handler
		cleanupsContext, ok := context.GetOk(r, "Cleanup")
		if ok {
			cleanups = cleanupsContext.([]func(h http.Handler) http.Handler)

			for i := 0; i < len(cleanups); i++ {
				cleanups[i](
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				).ServeHTTP(w, r)
			}
		}

		h.ServeHTTP(w, r)
	})
}
