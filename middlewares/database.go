package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
	"strconv"
)

func DatabaseMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

		dbOptions := &database.DatabaseInitOptions{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBname:   os.Getenv("DB_NAME"),
		}
		db, err := database.NewDatabase(dbOptions)
		if err != nil {
			utils.EncodeJSONError(w, "Failed to initialize database", http.StatusInternalServerError)
			return
		}

		context.Set(r, "DB", db)

		h.ServeHTTP(w, r)
	})
}
