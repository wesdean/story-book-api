package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
)

func ConfigMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configFilename := os.Getenv("CONFIG_FILENAME")
		if configFilename == "" {
			utils.EncodeJSONErrorWithLogging(r, w, "missing config filename", http.StatusInternalServerError)
			return
		}

		config, err := app_config.NewConfigFromFile(configFilename)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = config.Validate()
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}

		context.Set(r, "Config", config)

		h.ServeHTTP(w, r)
	})
}
