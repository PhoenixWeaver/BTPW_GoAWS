package types

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestNewUser tests the NewUser factory function
func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		input       RegisterUser
		expectError bool
	}{
		{
			name: "Valid user creation",
			input: RegisterUser{
				Username: "testuser",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Empty username",
			input: RegisterUser{
				Username: "",
				Password: "password123",
			},
			expectError: false, // NewUser doesn't validate input
		},
		{
			name: "Empty password",
			input: RegisterUser{
				Username: "testuser",
				Password: "",
			},
			expectError: false, // NewUser doesn't validate input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError {
				if user.Username != tt.input.Username {
					t.Errorf("Expected username %s, got %s", tt.input.Username, user.Username)
				}
				if user.PasswordHash == "" {
					t.Errorf("Password hash should not be empty")
				}
				if user.PasswordHash == tt.input.Password {
					t.Errorf("Password should be hashed, not stored as plain text")
				}
			}
		})
	}
}

// TestValidatePassword tests password validation
func TestValidatePassword(t *testing.T) {
	// Create a test user with hashed password
	registerUser := RegisterUser{
		Username: "testuser",
		Password: "password123",
	}
	user, err := NewUser(registerUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		expected       bool
	}{
		{
			name:           "Correct password",
			hashedPassword: user.PasswordHash,
			plainPassword:  "password123",
			expected:       true,
		},
		{
			name:           "Incorrect password",
			hashedPassword: user.PasswordHash,
			plainPassword:  "wrongpassword",
			expected:       false,
		},
		{
			name:           "Empty password",
			hashedPassword: user.PasswordHash,
			plainPassword:  "",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePassword(tt.hashedPassword, tt.plainPassword)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestCreateToken tests JWT token creation
func TestCreateToken(t *testing.T) {
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	token := CreateToken(user)
	if token == "" {
		t.Errorf("Token should not be empty")
	}

	// Parse the token to verify its structure
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(BT_JWT_SECRET), nil // Use the global JWT secret constant
	})

	if err != nil {
		t.Errorf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Errorf("Token should be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("Failed to extract claims")
	}

	if claims["user"] != "testuser" {
		t.Errorf("Expected user 'testuser', got %v", claims["user"])
	}

	// Check expiration
	exp, ok := claims["expires"].(float64)
	if !ok {
		t.Errorf("Expiration claim should be present")
	}

	// Token should expire in about 1 hour
	expectedExp := time.Now().Add(time.Hour).Unix()
	actualExp := int64(exp)

	// Allow 5 minutes tolerance
	if actualExp < expectedExp-300 || actualExp > expectedExp+300 {
		t.Errorf("Token expiration time is incorrect. Expected around %d, got %d", expectedExp, actualExp)
	}
}

// TestValidateToken tests JWT token validation
func TestValidateToken(t *testing.T) {
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	// Create a valid token
	token := CreateToken(user)

	tests := []struct {
		name        string
		tokenString string
		expectError bool
	}{
		{
			name:        "Valid token",
			tokenString: token,
			expectError: false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.here",
			expectError: true,
		},
		{
			name:        "Empty token",
			tokenString: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenString)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError {
				if claims["user"] != "testuser" {
					t.Errorf("Expected user 'testuser', got %v", claims["user"])
				}
			}
		})
	}
}
