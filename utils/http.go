package utils

import (
	"encoding/json"
	"fmt"
	"github.com/wesdean/story-book-api/logging"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, output interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		_ = EncodeJSONError(w, "Failed to encode output", http.StatusInternalServerError)
	}
}

func EncodeJSONWithStatus(w http.ResponseWriter, output interface{}, status int) {
	w.WriteHeader(status)
	EncodeJSON(w, output)
}

func EncodeJSONError(w http.ResponseWriter, errorMessage string, errorCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errorMessage), errorCode)
	}
	return err
}

func EncodeJSONErrorWithLogging(r *http.Request, w http.ResponseWriter, errorMessage string, errorCode int) {
	logger := logging.GetLoggerFromRequest(r)

	err := EncodeJSONError(w, errorMessage, errorCode)
	if err != nil {
		logging.Log(logger, logging.LOGLEVEL_ERROR, err.Error())
	}
	logging.Log(logger, logging.LOGLEVEL_ERROR, fmt.Sprintf("[%d] %s", errorCode, errorMessage))
}
