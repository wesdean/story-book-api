package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, output interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		EncodeJSONError(w, "Failed to encode output", http.StatusInternalServerError)
	}
}

func EncodeJSONError(w http.ResponseWriter, errorMessage string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errorMessage), errorCode)
	}
}
