package flow

import (
	"net/http"
)

// MiddlewareHandler is a custom middleware handler that satisfies the http.Handler interface.
type MiddlewareHandler struct {
	handler http.Handler
}

// ServeHTTP implements the http.Handler interface for MiddlewareHandler.
func (mh *MiddlewareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mh.handler.ServeHTTP(w, req)
}

// Middleware represents a flow middleware function.
type Middleware func(http.Handler) http.Handler

// ChainedMiddleware chains multiple middleware functions.
func ChainedMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
