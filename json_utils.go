package gorouter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONResponse sends a JSON response with the given data and status code.
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) // Set the status code once
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// JSONError writes an error message in JSON format to the client.
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	JSONResponse(w, map[string]string{"error": message}, statusCode)
}

// ParseJSONBody parses the JSON body of a request into the provided struct.
func ParseJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		JSONError(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return fmt.Errorf("content type must be application/json")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		JSONError(w, "Invalid JSON body", http.StatusBadRequest)
		return err
	}
	return nil
}
