package api

import (
	"BTgoAWSlambda/database"
	"BTgoAWSlambda/types"
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// TestApiHandler_RegisterUser tests the RegisterUser handler
func TestApiHandler_RegisterUser(t *testing.T) {
	// Create mock database
	mockStore := database.NewMockUserStore()
	handler := NewApiHandler(mockStore)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid registration",
			requestBody:    `{"username": "testuser", "password": "password123"}`,
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"username": "testuser"`, // Missing closing brace
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "Empty username",
			requestBody:    `{"username": "", "password": "password123"}`,
			expectedStatus: 400,
			expectError:    true, // API returns 400 status and also returns error
		},
		{
			name:           "Empty password",
			requestBody:    `{"username": "testuser", "password": ""}`,
			expectedStatus: 400,
			expectError:    true, // API returns 400 status and also returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Body: tt.requestBody,
			}

			response, err := handler.RegisterUser(request)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

// TestApiHandler_LoginUser tests the LoginUser handler
func TestApiHandler_LoginUser(t *testing.T) {
	// Create mock database with existing user
	mockStore := database.NewMockUserStore()
	handler := NewApiHandler(mockStore)

	// Create a test user with proper password hash
	registerUser := types.RegisterUser{
		Username: "testuser",
		Password: "password123",
	}
	user, err := types.NewUser(registerUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	mockStore.InsertUser(context.Background(), user)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid login",
			requestBody:    `{"username": "testuser", "password": "password123"}`,
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"username": "testuser"`, // Missing closing brace
			expectedStatus: 400,
			expectError:    true,
		},
		{
			name:           "Wrong password",
			requestBody:    `{"username": "testuser", "password": "wrongpassword"}`,
			expectedStatus: 401,
			expectError:    false,
		},
		{
			name:           "Non-existent user",
			requestBody:    `{"username": "nonexistent", "password": "password123"}`,
			expectedStatus: 500, // Database error
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Body: tt.requestBody,
			}

			response, err := handler.LoginUser(request)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

// TestApiHandler_RegisterUser_DuplicateUser tests duplicate user registration
func TestApiHandler_RegisterUser_DuplicateUser(t *testing.T) {
	// Create mock database with existing user
	mockStore := database.NewMockUserStore()
	handler := NewApiHandler(mockStore)

	// Create a test user in the mock database
	user := types.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	mockStore.InsertUser(context.Background(), user)

	// Try to register the same user again
	request := events.APIGatewayProxyRequest{
		Body: `{"username": "testuser", "password": "password123"}`,
	}

	response, err := handler.RegisterUser(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != 409 { // Conflict
		t.Errorf("Expected status 409 (Conflict), got %d", response.StatusCode)
	}
}

// TestApiHandler_LoginUser_InvalidPassword tests login with invalid password
func TestApiHandler_LoginUser_InvalidPassword(t *testing.T) {
	// Create mock database with existing user
	mockStore := database.NewMockUserStore()
	handler := NewApiHandler(mockStore)

	// Create a test user with proper password hash
	registerUser := types.RegisterUser{
		Username: "testuser",
		Password: "password123",
	}
	user, err := types.NewUser(registerUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	mockStore.InsertUser(context.Background(), user)

	// Try to login with wrong password
	request := events.APIGatewayProxyRequest{
		Body: `{"username": "testuser", "password": "wrongpassword"}`,
	}

	response, err := handler.LoginUser(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != 401 { // Unauthorized
		t.Errorf("Expected status 401 (Unauthorized), got %d", response.StatusCode)
	}
}

// TestApiHandler_RegisterUser_EmptyFields tests registration with empty fields
func TestApiHandler_RegisterUser_EmptyFields(t *testing.T) {
	mockStore := database.NewMockUserStore()
	handler := NewApiHandler(mockStore)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Empty username",
			requestBody:    `{"username": "", "password": "password123"}`,
			expectedStatus: 400,
			expectError:    true, // API returns 400 status and also returns error
		},
		{
			name:           "Empty password",
			requestBody:    `{"username": "testuser", "password": ""}`,
			expectedStatus: 400,
			expectError:    true, // API returns 400 status and also returns error
		},
		{
			name:           "Both empty",
			requestBody:    `{"username": "", "password": ""}`,
			expectedStatus: 400,
			expectError:    true, // API returns 400 status and also returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Body: tt.requestBody,
			}

			response, err := handler.RegisterUser(request)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}
