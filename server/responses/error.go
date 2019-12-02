package responses

import (
	"log"
	"net/http"
)

// Send an error response
func Error(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(`{"status": "error", "reason": "` + reason + `"}`)); err != nil {
		log.Printf("ERROR: failed to write response: %v\n", err)
	}
}
