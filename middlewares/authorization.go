package middlewares

import (
	"net/http"
)

func AuthorizationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//var authToken *AuthenticatedToken
		//authTokenContext, ok := context.GetOk(r, "AuthToken")
		//if ok {
		//	authToken = authTokenContext.(*AuthenticatedToken)
		//}

		h.ServeHTTP(w, r)
	})
}
