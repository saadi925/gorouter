// injector_test.go
package gorouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDependencyRegistry_Provide tests the Provide method of DependencyRegistry.
func TestDependencyRegistry_Provide(t *testing.T) {
	dr := NewDependencyRegistry()
	expected := "test dependency"

	dr.Provide("testKey", expected)

	actual, ok := dr.dependencies["testKey"]
	if !ok {
		t.Error("Dependency not registered")
	}

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

// TestDependencyRegistry_Middleware tests the Middleware method of DependencyRegistry.
func TestDependencyRegistry_Middleware(t *testing.T) {
	dr := NewDependencyRegistry()
	dr.Provide("testKey", "testValue")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := GetDependency(r.Context(), "testKey")
		if err != nil {
			t.Errorf("Error getting dependency from context: %v", err)
		}
		if val != "testValue" {
			t.Errorf("Expected %s, got %s", "testValue", val)
		}
	})

	mw := dr.Middleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	mw.ServeHTTP(recorder, req)
}

// TestGetDependencyFromContext tests the GetDependencyFromContext function.
func TestGetDependencyFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKey("testKey"), "testValue")

	val, err := GetDependency(ctx, "testKey")
	if err != nil {
		t.Errorf("Error getting dependency from context: %v", err)
	}
	if val != "testValue" {
		t.Errorf("Expected %s, got %s", "testValue", val)
	}
}

// Additional test cases can be added to cover edge cases and error handling scenarios.
