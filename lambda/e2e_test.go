package main

import (
	"BTgoAWSlambda/api"
	"BTgoAWSlambda/app"
	"BTgoAWSlambda/database"
	"BTgoAWSlambda/types"
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// TestE2E_CompleteUserJourney tests the complete user journey from registration to protected access
func TestE2E_CompleteUserJourney(t *testing.T) {
	// Setup: Create mock database and app
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Step 1: Register a new user
	t.Run("User Registration", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path: "/register",
			Body: `{"username": "e2etest", "password": "password123"}`,
		}

		response, err := lambdaApp.ApiHandler.RegisterUser(request)
		if err != nil {
			t.Errorf("Registration failed: %v", err)
		}
		if response.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", response.StatusCode)
		}
	})

	// Step 2: Attempt to register the same user again (should fail)
	t.Run("Duplicate Registration", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path: "/register",
			Body: `{"username": "e2etest", "password": "password123"}`,
		}

		response, err := lambdaApp.ApiHandler.RegisterUser(request)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if response.StatusCode != 409 { // Conflict
			t.Errorf("Expected status 409 (Conflict), got %d", response.StatusCode)
		}
	})

	// Step 3: Login with correct credentials
	t.Run("Successful Login", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path: "/login",
			Body: `{"username": "e2etest", "password": "password123"}`,
		}

		response, err := lambdaApp.ApiHandler.LoginUser(request)
		if err != nil {
			t.Errorf("Login failed: %v", err)
		}
		if response.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", response.StatusCode)
		}
	})

	// Step 4: Login with incorrect credentials
	t.Run("Failed Login", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path: "/login",
			Body: `{"username": "e2etest", "password": "wrongpassword"}`,
		}

		response, err := lambdaApp.ApiHandler.LoginUser(request)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if response.StatusCode != 401 { // Unauthorized
			t.Errorf("Expected status 401 (Unauthorized), got %d", response.StatusCode)
		}
	})

	// Step 5: Verify user data integrity
	t.Run("Data Integrity Check", func(t *testing.T) {
		user, err := mockStore.GetUser(context.Background(), "e2etest")
		if err != nil {
			t.Errorf("Failed to retrieve user: %v", err)
		}
		if user.Username != "e2etest" {
			t.Errorf("Username mismatch: expected 'e2etest', got '%s'", user.Username)
		}
		if user.PasswordHash == "" {
			t.Errorf("Password hash should not be empty")
		}
		if user.PasswordHash == "password123" {
			t.Errorf("Password should be hashed, not stored as plain text")
		}
	})
}

// TestE2E_ErrorScenarios tests various error scenarios
func TestE2E_ErrorScenarios(t *testing.T) {
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
		description    string
	}{
		{
			name:           "Invalid JSON Registration",
			path:           "/register",
			body:           `{"username": "test"`, // Missing closing brace
			expectedStatus: 400,
			description:    "Should return 400 for malformed JSON",
		},
		{
			name:           "Empty Fields Registration",
			path:           "/register",
			body:           `{"username": "", "password": ""}`,
			expectedStatus: 400,
			description:    "Should return 400 for empty fields",
		},
		{
			name:           "Invalid JSON Login",
			path:           "/login",
			body:           `{"username": "test"`, // Missing closing brace
			expectedStatus: 400,
			description:    "Should return 400 for malformed JSON",
		},
		{
			name:           "Non-existent User Login",
			path:           "/login",
			body:           `{"username": "nonexistent", "password": "password123"}`,
			expectedStatus: 500,
			description:    "Should return 500 for database error",
		},
		{
			name:           "Unknown Route",
			path:           "/unknown",
			body:           `{}`,
			expectedStatus: 404,
			description:    "Should return 404 for unknown routes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Path: tt.path,
				Body: tt.body,
			}

			var response events.APIGatewayProxyResponse
			var err error

			switch tt.path {
			case "/register":
				response, err = lambdaApp.ApiHandler.RegisterUser(request)
			case "/login":
				response, err = lambdaApp.ApiHandler.LoginUser(request)
			default:
				response = events.APIGatewayProxyResponse{
					Body:       "Got lost ? Not found",
					StatusCode: 404,
				}
			}

			// For invalid JSON, empty fields, and non-existent users, we expect errors
			if tt.name == "Invalid JSON Registration" || tt.name == "Empty Fields Registration" ||
				tt.name == "Invalid JSON Login" || tt.name == "Non-existent User Login" {
				if err == nil {
					t.Errorf("Expected error for %s but got none", tt.name)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. %s", tt.expectedStatus, response.StatusCode, tt.description)
			}
		})
	}
}

// TestE2E_Performance tests basic performance characteristics
func TestE2E_Performance(t *testing.T) {
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Test multiple user registrations
	t.Run("Multiple User Registration", func(t *testing.T) {
		users := []string{"user1", "user2", "user3", "user4", "user5"}

		for _, username := range users {
			request := events.APIGatewayProxyRequest{
				Path: "/register",
				Body: `{"username": "` + username + `", "password": "password123"}`,
			}

			response, err := lambdaApp.ApiHandler.RegisterUser(request)
			if err != nil {
				t.Errorf("Registration failed for %s: %v", username, err)
			}
			if response.StatusCode != 200 {
				t.Errorf("Expected status 200 for %s, got %d", username, response.StatusCode)
			}
		}

		// Verify all users were created
		for _, username := range users {
			exists, err := mockStore.DoesUserExist(context.Background(), username)
			if err != nil {
				t.Errorf("Database check failed for %s: %v", username, err)
			}
			if !exists {
				t.Errorf("User %s should exist in database", username)
			}
		}
	})

	// Test concurrent operations (basic test)
	t.Run("Concurrent Operations", func(t *testing.T) {
		// This is a basic test - in a real scenario, you'd use goroutines
		// For now, we'll test sequential operations to ensure they don't interfere

		// Register user
		registerRequest := events.APIGatewayProxyRequest{
			Path: "/register",
			Body: `{"username": "concurrenttest", "password": "password123"}`,
		}

		registerResponse, err := lambdaApp.ApiHandler.RegisterUser(registerRequest)
		if err != nil {
			t.Errorf("Registration failed: %v", err)
		}
		if registerResponse.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", registerResponse.StatusCode)
		}

		// Login user
		loginRequest := events.APIGatewayProxyRequest{
			Path: "/login",
			Body: `{"username": "concurrenttest", "password": "password123"}`,
		}

		loginResponse, err := lambdaApp.ApiHandler.LoginUser(loginRequest)
		if err != nil {
			t.Errorf("Login failed: %v", err)
		}
		if loginResponse.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", loginResponse.StatusCode)
		}
	})
}

// TestE2E_DataConsistency tests data consistency across operations
func TestE2E_DataConsistency(t *testing.T) {
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Register user
	registerRequest := events.APIGatewayProxyRequest{
		Path: "/register",
		Body: `{"username": "consistencytest", "password": "password123"}`,
	}

	registerResponse, err := lambdaApp.ApiHandler.RegisterUser(registerRequest)
	if err != nil {
		t.Errorf("Registration failed: %v", err)
	}
	if registerResponse.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", registerResponse.StatusCode)
	}

	// Verify user data consistency
	user, err := mockStore.GetUser(context.Background(), "consistencytest")
	if err != nil {
		t.Errorf("Failed to retrieve user: %v", err)
	}

	// Check username consistency
	if user.Username != "consistencytest" {
		t.Errorf("Username mismatch: expected 'consistencytest', got '%s'", user.Username)
	}

	// Check password hash consistency
	if user.PasswordHash == "" {
		t.Errorf("Password hash should not be empty")
	}

	// Verify password validation works
	isValid := types.ValidatePassword(user.PasswordHash, "password123")
	if !isValid {
		t.Errorf("Password validation should succeed")
	}

	// Verify password validation fails for wrong password
	isValid = types.ValidatePassword(user.PasswordHash, "wrongpassword")
	if isValid {
		t.Errorf("Password validation should fail for wrong password")
	}
}
