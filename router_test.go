// router_test.go
package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterWithRouteGroup(t *testing.T) {
	r := NewRouter()

	// Define a route group for /api
	apiGroup := r.Group("/api")
	apiGroup.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Test", "true")
			next.ServeHTTP(w, req)
		})
	})

	// // Add a GET route to /api/users
	// apiGroup.GET("/users", func(w http.ResponseWriter, req *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("GET /api/users"))
	// })

	// Add a POST route to /api/users
	apiGroup.POST("/users", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("POST /api/users"))
	})

	// Add a route group for /admin
	adminGroup := r.Group("/admin")
	adminGroup.Provide("AdminDependency", "admin_value")

	// Add a GET route to /admin/dashboard
	adminGroup.GET("/dashboard", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET /admin/dashboard"))
	})

	// Test GET /api/users route
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	// checkResponse(t, recorder, http.StatusOK, "GET /api/users", "X-Test", "true")

	// Test POST /api/users route
	req = httptest.NewRequest(http.MethodPost, "/api/users", nil)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	checkResponse(t, recorder, http.StatusCreated, "POST /api/users", "X-Test", "true")

	// Test GET /admin/dashboard route
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	checkResponse(t, recorder, http.StatusOK, "GET /admin/dashboard", "", "")

	// Test route group specific dependency
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	checkResponse(t, recorder, http.StatusOK, "GET /admin/dashboard", "", "")
}

func checkResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody string, expectedHeaderKey, expectedHeaderValue string) {
	t.Helper()

	if recorder.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, recorder.Code)
	}

	if expectedHeaderKey != "" {
		value := recorder.Header().Get(expectedHeaderKey)
		if value != expectedHeaderValue {
			t.Errorf("Expected header %s: %s, got %s", expectedHeaderKey, expectedHeaderValue, value)
		}
	}

	body := recorder.Body.String()
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}
