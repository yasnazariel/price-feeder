package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	// Test data
	testPayload := map[string]string{
		"message": "test response",
		"status":  "success",
	}

	// Create a mock HTTP response recorder
	w := httptest.NewRecorder()

	// Call the function
	RespondWithJSON(w, http.StatusOK, testPayload)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	// Check response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != testPayload["message"] {
		t.Errorf("Expected message %s, got %s", testPayload["message"], response["message"])
	}

	if response["status"] != testPayload["status"] {
		t.Errorf("Expected status %s, got %s", testPayload["status"], response["status"])
	}
}

func TestRespondWithError(t *testing.T) {
	// Test data
	testError := errors.New("test error message")
	expectedStatusCode := http.StatusBadRequest

	// Create a mock HTTP response recorder
	w := httptest.NewRecorder()

	// Call the function
	RespondWithError(w, expectedStatusCode, testError)

	// Check status code
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, got %d", expectedStatusCode, w.Code)
	}

	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	// Check response body
	var errorResponse ErrResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse.Error != testError.Error() {
		t.Errorf("Expected error message %s, got %s", testError.Error(), errorResponse.Error)
	}
}

func TestRespondWithError_DifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		error      error
	}{
		{"Internal Server Error", http.StatusInternalServerError, errors.New("internal error")},
		{"Not Found", http.StatusNotFound, errors.New("not found")},
		{"Unauthorized", http.StatusUnauthorized, errors.New("unauthorized")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondWithError(w, tc.statusCode, tc.error)

			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}

			var errorResponse ErrResponse
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			if err != nil {
				t.Errorf("Failed to unmarshal error response: %v", err)
			}

			if errorResponse.Error != tc.error.Error() {
				t.Errorf("Expected error message %s, got %s", tc.error.Error(), errorResponse.Error)
			}
		})
	}
}
