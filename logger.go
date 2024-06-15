package flow

import (
	"log"
	"net/http"
)

// RequestLogger logs incoming requests.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Request: %s %s", req.Method, req.URL.Path)
		next.ServeHTTP(w, req)
	})
}
