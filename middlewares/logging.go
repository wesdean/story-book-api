package middlewares

import (
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
)

func LoggingMiddleware(h http.Handler) http.Handler {
	hFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var config *app_config.Config
		configVal, ok := context.GetOk(r, "Config")
		if ok {
			config = configVal.(*app_config.Config)
		}

		if config == nil {
			utils.EncodeJSONErrorWithLogging(r, w, "missing application config", http.StatusInternalServerError)
			return
		}

		logConfig, err := config.GetLogger("Config.API.Logger")
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger, err := logging.NewLogger(logConfig)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}

		context.Set(r, "Logger", logger)

		h.ServeHTTP(w, r)
	})
	return LoggingCleanupMiddleware(hFunc)
}