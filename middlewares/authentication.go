package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AuthenticatedToken struct {
	UserId    int
	Timestamp int64
}

func AuthenticationtMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		secret := []byte(os.Getenv("AUTH_SECRET"))

		authTimeout, err := strconv.Atoi(os.Getenv("AUTH_TIMEOUT"))
		if err != nil {
			utils.EncodeJSONError(w, "Invaild authentication timeout", http.StatusInternalServerError)
			return
		}

		if tokenString == "" {
			utils.EncodeJSONError(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWTToken(tokenString, secret)
		if err != nil {
			utils.EncodeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tsFloat, ok := claims["timestamp"].(float64)
		if !ok {
			utils.EncodeJSONError(w, "Invalid authentication timestamp", http.StatusUnauthorized)
			return
		}
		timestamp := int64(tsFloat)

		if (time.Now().Unix() - timestamp) > int64(authTimeout) {
			utils.EncodeJSONError(w, "Authentication has expired", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["user_id"].(int)
		authToken := AuthenticatedToken{
			UserId:    userId,
			Timestamp: timestamp,
		}

		context.Set(r, "AuthToken", authToken)

		h.ServeHTTP(w, r)
	})
}
