package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// Send a generic success response
func Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status": "success"}`)); err != nil {
		log.Printf("ERROR: failed to write response: %v\n", err)
	}
}

// Send a success response with some data
func SuccessWithData(w http.ResponseWriter, data interface{}) {
	// Encode body
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: failed to encode response data: %v\n", err)
		return
	}

	// Create the response
	response := append([]byte(`{"status": "success", "data": `), encoded...)
	response = append(response, []byte(`}`)...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Printf("ERROR: failed to write response: %v\n", err)
	}
}
