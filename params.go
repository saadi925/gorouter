package flow

import (
	"net/http"
	"strconv"
	"strings"
)

type ContextKey string

// GetStringSlice returns the slice of string values for the specified parameter.
func (p Params) GetStringSlice(key string) []string {
	val, ok := p[key]
	if !ok {
		return nil
	}
	return []string{val}
}

// GetIntSlice returns the slice of integer values for the specified parameter.
// If any parameter value is not a valid integer, it is skipped.
func (p Params) GetIntSlice(key string) []int {
	vals, ok := p[key]
	if !ok {
		return nil
	}

	var intSlice []int
	for _, val := range strings.Split(vals, ",") {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			intSlice = append(intSlice, intVal)
		}
	}
	return intSlice
}

// ParamsContextKey is the context key for request parameters.
const ParamsContextKey ContextKey = "params"

func GetParams(req *http.Request) Params {
	params := req.Context().Value(ParamsContextKey).(Params)
	return params
}

// Params represents the parameters extracted from the URL path.
type Params map[string]string

// Get returns the value of the specified parameter.
func (p Params) Get(key string) string {
	return p[key]
}

// GetInt returns the integer value of the specified parameter.
// If the parameter value is not a valid integer, it returns 0.
func (p Params) GetInt(key string) int {
	val, err := strconv.Atoi(p[key])
	if err != nil {
		return 0
	}
	return val
}

// parseParams extracts parameters from the request URL path.
func parseParams(req *http.Request) Params {
	params := make(Params)
	path := req.URL.Path

	parts := strings.Split(path, "/")
	for i := 1; i < len(parts); i += 2 { // Assuming parameters are in pairs
		if i+1 < len(parts) {
			key := strings.TrimPrefix(parts[i], ":") // Remove ":" from parameter name
			params[key] = parts[i+1]
		}
	}
	return params
}

// PathParam retrieves the value of a path parameter from the request context.
func PathParam(req *http.Request, key string) string {
	params := req.Context().Value(ParamsContextKey).(Params)
	return params.Get(key)
}
