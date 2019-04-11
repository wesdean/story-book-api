package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
)

func DatabaseCleanupMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerFromRequest(r)
		dbValue, ok := context.GetOk(r, "DB")
		if ok {
			db := dbValue.(*database.Database)
			err := db.Rollback()
			if err != nil {
				utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
				return
			}
			logging.Log(logger, logging.LOGLEVEL_INFO, "database closed")
		}

		h.ServeHTTP(w, r)
	})
}
