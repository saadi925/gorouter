package security

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
)

// CSRFTokenLength defines the length of CSRF tokens
const CSRFTokenLength = 32

// csrfTokenMap stores CSRF tokens
var csrfTokenMap = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

// GenerateCSRFToken generates a new CSRF token and stores it
func GenerateCSRFToken() (string, error) {
	token := make([]byte, CSRFTokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	encodedToken := base64.StdEncoding.EncodeToString(token)
	csrfTokenMap.Lock()
	csrfTokenMap.m[encodedToken] = true
	csrfTokenMap.Unlock()
	return encodedToken, nil
}

// ValidateCSRFToken validates the CSRF token
func ValidateCSRFToken(token string) bool {
	csrfTokenMap.RLock()
	defer csrfTokenMap.RUnlock()
	_, ok := csrfTokenMap.m[token]
	return ok
}

// CSRFMiddleware adds CSRF protection middleware.
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost || req.Method == http.MethodPut ||
			req.Method == http.MethodDelete || req.Method == http.MethodPatch {
			token := req.FormValue("csrf_token")
			if !ValidateCSRFToken(token) {
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}

// CSRFTokenMiddleware generates a CSRF token and sets it in the response header.
func CSRFTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := GenerateCSRFToken()
		if err != nil {
			http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
			return
		}
		// Set the token in a response header (or cookie, depending on your needs)
		w.Header().Set("X-CSRF-Token", token)
		next.ServeHTTP(w, req)
	})
}
