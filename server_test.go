package main

import (
	"forum/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	// Set the flag to indicate test mode.
	IsTest = true
}

// TestRoutes verifies that each route returns the expected HTTP status code.
// Note: For simplicity, this test assumes that the handlers return a status code of 200 (OK)
// when invoked with the provided (mock) data. In a real application, you may need to inject
// mocks or stubs to simulate dependencies such as the database, session data, etc.
func TestRoutes(t *testing.T) {
	// Register the routes by calling setupRoutes.
	_, sessionToken, _ := testInit()
	routes.SetupRoutes()

	testSessionCookie := &http.Cookie{
		Name:  "session_token", // Adjust the name to match your session cookie name.
		Value: sessionToken,    // You might need to set a value that your controllers expect.
		Path:  "/",
	}

	// Define a list of test cases. For POST methods, we use mock form data in the body.
	tests := []struct {
		name       string
		method     string
		target     string
		body       string // Form data for POST requests
		wantStatus int
		hasCookie  bool
	}{
		{"MainPage GET", http.MethodGet, "/", "", http.StatusOK, false},
		{"AuthHandler GET", http.MethodGet, "/auth/", "", http.StatusOK, false},
		{"First post GET", http.MethodGet, "/post/f9edb8d6-c739-4d6f-aaa4-9b298f2e1552", "", http.StatusOK, false},

		{"Register GET", http.MethodPost, "/register", "username=test2&email=test2@test2.com&password=test2&user_type=test_user", http.StatusFound, false},
		{"Login GET", http.MethodPost, "/login", "username=test2&password=test22", http.StatusOK, false},   //failed login
		{"Login GET", http.MethodPost, "/login", "username=test2&password=test2", http.StatusFound, false}, //login with user test2 that just registered
		// // For POST routes, we supply mock data as the body.
		{"SubmitPost POST", http.MethodPost, "/submitPost", "title=test&description=example&categories=1", http.StatusFound, true}, //submitPost with user test
	}

	// Iterate through each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			// Create a new HTTP request with the appropriate method and body.
			if tt.method == http.MethodPost {
				// Use mock data for POST requests.
				req = httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req = httptest.NewRequest(tt.method, tt.target, nil)
			}
			if tt.hasCookie {
				req.AddCookie(testSessionCookie)
			}

			// Create a ResponseRecorder to record the HTTP response.
			recorder := httptest.NewRecorder()

			// Use the default HTTP ServeMux where the routes were registered.
			http.DefaultServeMux.ServeHTTP(recorder, req)

			// Check if the HTTP status code is as expected.
			if recorder.Code != tt.wantStatus {
				t.Errorf("For route %s with method %s, expected status %d, but got %d", tt.target, tt.method, tt.wantStatus, recorder.Code)
			}
		})
	}

	finishTest()
}
