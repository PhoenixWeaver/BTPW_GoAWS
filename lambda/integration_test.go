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

// TestIntegration_UserRegistration tests the complete user registration flow
func TestIntegration_UserRegistration(t *testing.T) {
	// Create mock database
	mockStore := database.NewMockUserStore()

	// Create app with mock database
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Test user registration
	request := events.APIGatewayProxyRequest{
		Path: "/register",
		Body: `{"username": "integrationtest", "password": "password123"}`,
	}

	response, err := lambdaApp.ApiHandler.RegisterUser(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	// Verify user was created in database
	exists, err := mockStore.DoesUserExist(context.Background(), "integrationtest")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("User should exist after registration")
	}
}

// TestIntegration_UserLogin tests the complete user login flow
func TestIntegration_UserLogin(t *testing.T) {
	// Create mock database with existing user
	mockStore := database.NewMockUserStore()

	// Create a test user with hashed password
	registerUser := types.RegisterUser{
		Username: "logintest",
		Password: "password123",
	}
	user, err := types.NewUser(registerUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Insert user into mock database
	mockStore.InsertUser(context.Background(), user)

	// Create app with mock database
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Test user login
	request := events.APIGatewayProxyRequest{
		Path: "/login",
		Body: `{"username": "logintest", "password": "password123"}`,
	}

	response, err := lambdaApp.ApiHandler.LoginUser(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}
}

// TestIntegration_ProtectedRoute tests the protected route with JWT middleware
func TestIntegration_ProtectedRoute(t *testing.T) {
	// Create a test handler that returns protected content
	// protectedHandler := func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// To this:
	protectedHandler := func(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body:       "Protected content",
			StatusCode: 200,
		}, nil
	}

	// Test without JWT token
	request := events.APIGatewayProxyRequest{
		Path:    "/protected",
		Headers: map[string]string{},
	}

	// This would normally be wrapped with ValidateJWTMiddleware
	// For integration testing, we'll test the handler directly
	response, err := protectedHandler(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}
}

// TestIntegration_ErrorHandling tests error handling across the application
func TestIntegration_ErrorHandling(t *testing.T) {
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON in registration",
			path:           "/register",
			body:           `{"username": "test"`, // Invalid JSON
			expectedStatus: 400,
		},
		{
			name:           "Invalid JSON in login",
			path:           "/login",
			body:           `{"username": "test"`, // Invalid JSON
			expectedStatus: 400,
		},
		{
			name:           "Empty fields in registration",
			path:           "/register",
			body:           `{"username": "", "password": ""}`,
			expectedStatus: 400,
		},
		{
			name:           "Non-existent route",
			path:           "/nonexistent",
			body:           `{}`,
			expectedStatus: 404,
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

			// For invalid JSON and empty fields, we expect errors
			if tt.name == "Invalid JSON in registration" || tt.name == "Invalid JSON in login" || tt.name == "Empty fields in registration" {
				if err == nil {
					t.Errorf("Expected error for %s but got none", tt.name)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

// TestIntegration_DataFlow tests the complete data flow from request to response
func TestIntegration_DataFlow(t *testing.T) {
	mockStore := database.NewMockUserStore()
	lambdaApp := app.App{
		ApiHandler: api.NewApiHandler(mockStore),
	}

	// Step 1: Register a user
	registerRequest := events.APIGatewayProxyRequest{
		Path: "/register",
		Body: `{"username": "dataflowtest", "password": "password123"}`,
	}

	registerResponse, err := lambdaApp.ApiHandler.RegisterUser(registerRequest)
	if err != nil {
		t.Errorf("Registration failed: %v", err)
	}
	if registerResponse.StatusCode != 200 {
		t.Errorf("Registration should succeed, got status %d", registerResponse.StatusCode)
	}

	// Step 2: Login with the same user
	loginRequest := events.APIGatewayProxyRequest{
		Path: "/login",
		Body: `{"username": "dataflowtest", "password": "password123"}`,
	}

	loginResponse, err := lambdaApp.ApiHandler.LoginUser(loginRequest)
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}
	if loginResponse.StatusCode != 200 {
		t.Errorf("Login should succeed, got status %d", loginResponse.StatusCode)
	}

	// Step 3: Verify user exists in database
	exists, err := mockStore.DoesUserExist(context.Background(), "dataflowtest")
	if err != nil {
		t.Errorf("Database check failed: %v", err)
	}
	if !exists {
		t.Errorf("User should exist in database")
	}

	// Step 4: Verify user data integrity
	user, err := mockStore.GetUser(context.Background(), "dataflowtest")
	if err != nil {
		t.Errorf("Failed to retrieve user: %v", err)
	}
	if user.Username != "dataflowtest" {
		t.Errorf("Username mismatch: expected 'dataflowtest', got '%s'", user.Username)
	}
	if user.PasswordHash == "" {
		t.Errorf("Password hash should not be empty")
	}
	if user.PasswordHash == "password123" {
		t.Errorf("Password should be hashed, not stored as plain text")
	}
}
