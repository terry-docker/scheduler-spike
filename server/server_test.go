// httpserver/server_test.go
package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	// Create a handler with a mock isRunning function
	handler := StatusHandler(func() bool {
		return true // Simulate that the scheduler is running
	})

	// Use httptest to create a request and response recorder
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Service is running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test the "not running" scenario
	handler = StatusHandler(func() bool {
		return false // Simulate that the scheduler is not running
	})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the response for the "not running" case
	expected = "Service is not running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
