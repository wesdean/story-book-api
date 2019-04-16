package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func AuthenticationtMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		secret := []byte(os.Getenv("AUTH_SECRET"))

		authTimeout, err := strconv.Atoi(os.Getenv("AUTH_TIMEOUT"))
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, "invaild authentication timeout", http.StatusInternalServerError)
			return
		}

		if tokenString == "" {
			utils.EncodeJSONErrorWithLogging(r, w, "invalid authentication token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWTToken(tokenString, secret)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusUnauthorized)
			return
		}

		tsFloat, ok := claims["timestamp"].(float64)
		if !ok {
			utils.EncodeJSONErrorWithLogging(r, w, "invalid authentication timestamp", http.StatusUnauthorized)
			return
		}
		timestamp := int64(tsFloat)

		if (time.Now().Unix() - timestamp) > int64(authTimeout) {
			utils.EncodeJSONErrorWithLogging(r, w, "authentication has expired", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["user_id"].(float64)
		if !ok {
			utils.EncodeJSONErrorWithLogging(r, w, "invalid user id in token", http.StatusBadRequest)
			return
		}
		stores, err := models.GetStoresFromRequest(r)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := stores.UserStore.GetUser(models.NewUserQueryOptions().Id(int(userId)))
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user.Disabled {
			utils.EncodeJSONErrorWithLogging(r, w, "user is disabled", http.StatusUnauthorized)
			return
		}

		if user.Archived {
			utils.EncodeJSONErrorWithLogging(r, w, "user is archived", http.StatusUnauthorized)
			return
		}

		authenticatedUser := &models.AuthenticatedUser{
			User:      user,
			Timestamp: int(timestamp),
		}

		context.Set(r, "AuthenticatedUser", authenticatedUser)

		h.ServeHTTP(w, r)
	})
}

func GetAuthenticatedUserFromRequest(r *http.Request) *models.AuthenticatedUser {
	authUserContext, ok := context.GetOk(r, "AuthenticatedUser")
	if ok {
		authUser, ok := authUserContext.(*models.AuthenticatedUser)
		if ok {
			return authUser
		}
	}
	return nil
}
