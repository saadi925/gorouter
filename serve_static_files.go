package gorouter

import (
	"net/http"
	"strings"
)

// StaticFileServer returns a handler that serves static files from the given directory.
func StaticFileServer(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Ensure the request path is safe to prevent directory traversal attacks
		if strings.Contains(req.URL.Path, "..") {
			http.NotFound(w, req)
			return
		}
		fs.ServeHTTP(w, req)
	})
}
