package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
	"strconv"
)

func DatabaseMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerFromRequest(r)
		port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
		dbOptions := &database.DatabaseInitOptions{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBname:   os.Getenv("DB_NAME"),
		}

		logging.Log(logger, logging.LOGLEVEL_INFO, "initializing database connection")
		db, err := database.NewDatabase(dbOptions)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w,
				"failed to initialize database connection",
				http.StatusInternalServerError)
			return
		}
		logging.Log(logger, logging.LOGLEVEL_INFO, "database initialized")

		err = db.Begin()
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
		logging.Log(logger, logging.LOGLEVEL_INFO, "database opened")

		logging.Log(logger, logging.LOGLEVEL_INFO, "stores created")
		stores := models.NewStores(db, logger)

		context.Set(r, "DB", db)
		context.Set(r, "Stores", stores)

		AppendCleanups(r, DatabaseCleanupMiddleware)
		h.ServeHTTP(w, r)
	})
}
