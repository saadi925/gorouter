package security

import (
	"net/http"
	"strings"
)

// CORSOptions represents the available options for CORS configuration.
type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// CORSMiddleware creates a middleware to handle CORS based on provided options.
func CORSMiddleware(options CORSOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if origin == "" {
				// No origin header present, so skip CORS checks
				next.ServeHTTP(w, req)
				return
			}

			if isOriginAllowed(options.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ","))

				if req.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.WriteHeader(http.StatusNoContent)
					return
				}
			} else {
				// Origin not allowed
				http.Error(w, "CORS origin not allowed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// isOriginAllowed checks if the origin is allowed based on the allowed origins list.
func isOriginAllowed(allowedOrigins []string, origin string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
