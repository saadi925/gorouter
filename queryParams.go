// flow/query_params.go
package flow

import (
	"net/http"
	"strconv"
)

// QueryParams represents query parameters in a map format.
type QueryParams map[string]string

// Get returns the value of the specified query parameter.
func (qp QueryParams) Get(key string) string {
	return qp[key]
}

// GetInt returns the integer value of the specified query parameter.
// If the parameter value is not a valid integer, it returns 0.
func (qp QueryParams) GetInt(key string) int {

	val, err := strconv.Atoi(qp[key])
	if err != nil {
		return 0
	}
	return val
}

// ParseQueryParams parses the query parameters from the request URL.
func ParseQueryParams(req *http.Request) QueryParams {
	params := make(QueryParams)
	query := req.URL.Query()

	for key, values := range query {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}
