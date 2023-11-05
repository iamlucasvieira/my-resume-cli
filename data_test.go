package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock server response
var mockResponse = CV{
	Info: []Info{
		{
			Name:  "Email",
			Value: "lucas6eng@gmail.com",
			Alias: []string{"email", "e-mail"},
		},
	},
}

// TestRequest tests the request method of the api struct
func TestRequest(t *testing.T) {
	// Create a mock serer
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(mockResponse)
		_, err := w.Write(resp)
		if err != nil {
			http.Error(w, "Failed to write response of mock server", http.StatusInternalServerError)
		}
	}))

	// Close the server when test finishes
	defer mockServer.Close()

	// Tests the API
	a := api{}
	var cv CV

	err := a.request(mockServer.URL, &cv) // Use server.URL as the URL to request
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	// Add assertions based on the expected content of cv
	if len(cv.Info) != 1 || cv.Info[0].Name != "Email" {
		t.Errorf("Unexpected CV data: %+v", cv)
	}
}
