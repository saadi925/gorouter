package flow

import (
	"context"
	"net/http"
	"strings"
)

// Router represents the flow router for managing routes.
type Router struct {
	*http.ServeMux
	middleware         []Middleware
	routes             map[string]map[string]routeHandler
	globalDependencies *DependencyRegistry
}

// routeHandler holds the handler and its specific dependencies
type routeHandler struct {
	handler            http.HandlerFunc
	dependencyRegistry *DependencyRegistry
}

// NewRouter creates a new instance of the flow router.
func NewRouter() *Router {
	return &Router{
		middleware:         []Middleware{},
		routes:             make(map[string]map[string]routeHandler),
		ServeMux:           http.NewServeMux(),
		globalDependencies: NewDependencyRegistry(), // Initialize global DependencyRegistry
	}
}

// Group creates a new route group with a common prefix.
func (r *Router) Group(prefix string, middleware ...Middleware) *RouteGroup {
	return &RouteGroup{
		prefix:     prefix,
		middleware: middleware,
		router:     r,
	}
}

// Use applies middleware to the router.
func (r *Router) Use(middleware ...Middleware) {
	r.middleware = append(r.middleware, middleware...)
}

// AddRoute adds a new route to the router.
func (r *Router) AddRoute(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	finalHandler := ApplyMiddleware(handler, append(r.middleware, middleware...)...)

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]routeHandler)
	}
	r.routes[method][path] = routeHandler{
		handler:            finalHandler.ServeHTTP,
		dependencyRegistry: r.globalDependencies, // Use global DependencyRegistry
	}
	r.ServeMux.Handle(path, finalHandler)
}

// AddRouteWithDependencies adds a new route to the router with route-specific dependencies.
func (r *Router) AddRouteWithDependencies(method, path string, handler http.HandlerFunc, dependencies *DependencyRegistry, middleware ...Middleware) {
	finalHandler := ApplyMiddleware(handler, append(r.middleware, middleware...)...)

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]routeHandler)
	}
	r.routes[method][path] = routeHandler{
		handler:            finalHandler.ServeHTTP,
		dependencyRegistry: dependencies, // Use provided route-specific DependencyRegistry
	}
	r.ServeMux.Handle(path, finalHandler)
}

// ServeHTTP handles incoming HTTP requests and dispatches them to the appropriate handlers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := parseParams(req)
	ctx := req.Context()
	ctx = context.WithValue(ctx, ParamsContextKey, params)
	req = req.WithContext(ctx)

	finalHandler := ApplyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Try to match the exact path first
		if rh, ok := r.routes[req.Method][req.URL.Path]; ok {
			mergedHandler := r.mergeHandlersWithDependencies(rh.handler, rh.dependencyRegistry)
			mergedHandler(w, req)
			return
		}

		// Try to match with parameters in the path
		for path, rh := range r.routes[req.Method] {
			if matched, parsedParams := matchPathWithParams(path, req.URL.Path); matched {
				req = req.WithContext(context.WithValue(req.Context(), ParamsContextKey, parsedParams))
				mergedHandler := r.mergeHandlersWithDependencies(rh.handler, rh.dependencyRegistry)
				mergedHandler(w, req)
				return
			}
		}

		http.NotFound(w, req)
	}), r.middleware...)

	finalHandler.ServeHTTP(w, req)
}

// mergeHandlersWithDependencies merges global and route-specific dependencies
func (r *Router) mergeHandlersWithDependencies(handler http.HandlerFunc, routeDependencies *DependencyRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		// Apply global dependencies
		for key, value := range r.globalDependencies.dependencies {
			ctx = context.WithValue(ctx, ContextKey(key), value)
		}
		// Apply route-specific dependencies
		if routeDependencies != nil {
			for key, value := range routeDependencies.dependencies {
				ctx = context.WithValue(ctx, ContextKey(key), value)
			}
		}
		req = req.WithContext(ctx)
		handler(w, req)
	}
}

// matchPathWithParams matches the path with parameters against the request path.
func matchPathWithParams(routePath, requestPath string) (bool, Params) {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return false, nil
	}

	params := make(Params)
	for i := 0; i < len(routeParts); i++ {
		routePart := routeParts[i]
		requestPart := requestParts[i]

		if strings.HasPrefix(routePart, ":") {
			key := strings.TrimPrefix(routePart, ":")
			params[key] = requestPart
		} else if routePart != requestPart {
			return false, nil
		}
	}

	return true, params
}

// ApplyMiddleware applies middleware to a handler.
func ApplyMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, mw := range middleware {
		handler = mw(handler)
	}
	return handler
}
