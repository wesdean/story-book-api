package middlewares

import (
	"github.com/wesdean/story-book-api/logging"
	"net/http"
)

func LoggingCleanupMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerFromRequest(r)
		if logger != nil {
			logging.CloseLogger(logger)
		}

		h.ServeHTTP(w, r)
	})
}
