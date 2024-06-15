package flow

import (
	"context"
	"errors"
	"net/http"
)

// DependencyRegistry is responsible for managing application dependencies.
type DependencyRegistry struct {
	dependencies map[string]interface{}
}

// NewDependencyRegistry creates a new instance of DependencyRegistry.
func NewDependencyRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		dependencies: make(map[string]interface{}),
	}
}

// Provide registers a dependency with a key.
func (dr *DependencyRegistry) Provide(key string, dependency interface{}) {
	dr.dependencies[key] = dependency
}

// Middleware injects registered dependencies into the request context.
func (dr *DependencyRegistry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for key, value := range dr.dependencies {
			ctx = context.WithValue(ctx, ContextKey(key), value)
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// GetDependencyFromContext retrieves a dependency from the request context.
func GetDependencyFromContext(ctx context.Context, key string) (interface{}, error) {
	val := ctx.Value(ContextKey(key))
	if val == nil {
		return nil, errors.New("dependency not found in context")
	}
	return val, nil
}
