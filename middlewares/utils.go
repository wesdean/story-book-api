package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
)

//region Common Responses
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	utils.EncodeJSONErrorWithLogging(r, w, "Access denied", http.StatusUnauthorized)
}

//endregion

//region Authentication Helpers
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

func IsUserAuthenticated(r *http.Request) bool {
	authUser := GetAuthenticatedUserFromRequest(r)
	return authUser != nil
}

//endregion
