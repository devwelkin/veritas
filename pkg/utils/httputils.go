package utils // or appropriate package name for the project

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON writes any data as JSON with the specified status code.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// convert payload
	dat, err := json.Marshal(payload)
	if err != nil {
		// if our data (payload) cannot be converted to JSON, this is a server error
		log.Printf("failed to marshal json response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// write JSON data to response
	w.Write(dat)
}

// RespondWithError returns a standard JSON error message.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	// returning errors in the same format ensures consistency
	errorPayload := map[string]string{"error": message}
	RespondWithJSON(w, code, errorPayload)
}
