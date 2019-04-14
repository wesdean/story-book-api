package middlewares

import (
	"github.com/wesdean/story-book-api/utils"
	"net/http"
)

func AuthorizationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := GetAuthenticatedUserFromRequest(r)
		requestID := r.Method + " " + r.URL.Path

		// Default is to reject authorization
		authorized := false
		rejectionMessage := "User not authorized for this action"

		switch requestID {
		// Does not require authentication
		case "GET /":
			authorized = true

		// Requires authentication only
		case "GET /user_roles":
			if authenticatedUser != nil {
				authorized = true
			}
		}

		if authorized {
			h.ServeHTTP(w, r)
		} else {
			utils.EncodeJSONErrorWithLogging(r, w, rejectionMessage, http.StatusUnauthorized)
			return
		}

	})
}
