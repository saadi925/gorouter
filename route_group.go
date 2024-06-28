package gorouter

import (
	"net/http"
)

// RouteGroup represents a group of routes.
type RouteGroup struct {
	prefix             string
	middleware         []Middleware
	router             *Router
	dependencyRegistry *DependencyRegistry
}

// Use adds middleware to the route group.
func (rg *RouteGroup) Use(middleware ...Middleware) {
	rg.middleware = append(rg.middleware, middleware...)
}

// Provide registers a route-specific dependency.
func (rg *RouteGroup) Provide(key string, dependency interface{}) {
	if rg.dependencyRegistry == nil {
		rg.dependencyRegistry = NewDependencyRegistry()
	}
	rg.dependencyRegistry.Provide(key, dependency)
}

// Group creates a new route group with a common prefix.
func (rg *RouteGroup) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:             rg.prefix + prefix,
		middleware:         rg.middleware,
		router:             rg.router,
		dependencyRegistry: rg.dependencyRegistry,
	}
}

// GET adds a GET route to the route group.
func (rg *RouteGroup) GET(path string, handler http.HandlerFunc, middleware ...Middleware) {
	rg.router.AddRouteWithDependencies(http.MethodGet, rg.prefix+path, handler, rg.dependencyRegistry, append(rg.middleware, middleware...)...)
}

// POST adds a POST route to the route group.
func (rg *RouteGroup) POST(path string, handler http.HandlerFunc, middleware ...Middleware) {
	rg.router.AddRouteWithDependencies(http.MethodPost, rg.prefix+path, handler, rg.dependencyRegistry, append(rg.middleware, middleware...)...)
}
