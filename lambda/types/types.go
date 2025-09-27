/*
===============================================================================
AWS LAMBDA FUNCTION - TYPES PACKAGE
===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This package contains all data structures, types, and utility functions
             for user authentication, password hashing, and JWT token management
             in our AWS Lambda function.

LEARNING OBJECTIVES:
- Data structure design for Lambda functions
- Password security with bcrypt hashing
- JWT token creation and validation
- Error handling patterns in Go
- Type safety and validation

TOPICS COVERED:
- Struct definitions with JSON tags
- Password hashing with bcrypt
- JWT token creation and claims
- Error handling and validation
- Type conversion and serialization

FEATURES:
- User registration data structure
- Secure password hashing
- JWT token generation
- Password validation
- Type-safe data handling

ARCHITECTURE:
- RegisterUser: Input data for user registration
- User: Internal user data structure with hashed password
- NewUser: Factory function to create users with hashed passwords
- ValidatePassword: Password verification function
- CreateToken: JWT token generation function

FUNCTION RELATIONSHIPS:
1. RegisterUser → NewUser → User (Registration flow)
2. User + plain password → ValidatePassword → bool (Login flow)
3. User → CreateToken → JWT string (Token generation)

*/

package types

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// BT_JWT_SECRET is the secret key used for signing and validating JWT tokens.
// TODO: In production, move this to environment variables for security.
// lambda/types/types.go: const BT_JWT_SECRET = "{@BTPhoenix}"
// Used in: CreateToken() and ValidateToken() functions
// Used in: lambda/middleware/middleware.go via types.BT_JWT_SECRET
// Used in: lambda/types/types_test.go for testing
// Future Implementation (Environment Variables):
// From the documentation files, I found references to:
// jwtSecret := os.Getenv("JWT_SECRET")

const BT_JWT_SECRET = "{@BTPhoenix}"

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// RegisterUser represents the input data structure for user registration.
// This is what clients send when creating a new user account.
//
// JSON Tags: Used for automatic serialization/deserialization
// - "username": Maps to Username field in JSON
// - "password": Maps to Password field in JSON
//
// Usage: This struct is used in the /register endpoint
type RegisterUser struct {
	Username string `json:"username"` // User's chosen username (required)
	Password string `json:"password"` // User's plain text password (will be hashed)
}

// User represents the internal user data structure stored in the database.
// This contains the hashed password for security.
//
// Security Note: PasswordHash should NEVER be returned to clients
// Only use this struct internally for database operations
//
// Usage: This struct is used for database storage and internal operations
type User struct {
	Username     string `json:"username" dynamodbav:"username"` // User's username (same as RegisterUser)
	PasswordHash string `json:"password" dynamodbav:"password"` // Hashed password (NOT plain text)
}

// ============================================================================
// FACTORY FUNCTIONS
// ============================================================================

// NewUser creates a new User struct from RegisterUser input.
// This function handles password hashing using bcrypt for security.
//
// Parameters:
//   - registerUser: The input data from user registration
//
// Returns:
//   - User: The created user with hashed password
//   - error: Any error that occurred during password hashing
//
// Security Features:
// - Uses bcrypt with cost factor 10 (good balance of security/performance)
// - Automatically handles salt generation
// - Returns error if hashing fails
//
// Usage: Called during user registration process
func NewUser(registerUser RegisterUser) (User, error) {
	// Hash the password with bcrypt
	// Cost factor 10 = 2^10 = 1024 iterations (good security level)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		// Return empty User and the error if hashing fails
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create and return the User struct with hashed password
	return User{
		Username:     registerUser.Username,  // Username from RegisterUser struct
		PasswordHash: string(hashedPassword), // Convert []byte to string
	}, nil // return user, nil if no error
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// ValidatePassword compares a plain text password with a hashed password.
// This function is used during user login to verify credentials.
//
// Parameters:
//   - hashedPassword: The stored hashed password from database
//   - plainTextPassword: The password provided by user during login
//
// Returns:
//   - bool: true if passwords match, false otherwise
//
// Security Features:
// - Uses bcrypt.CompareHashAndPassword for secure comparison
// - Timing-safe comparison (prevents timing attacks)
// - Automatically handles salt extraction from hash
//
// Usage: Called during user login process
func ValidatePassword(hashedPassword, plainTextPassword string) bool {
	// Compare the plain text password with the stored hash
	// bcrypt automatically extracts salt from the hash
	// CompareHashAndPassword compares a bcrypt hash with a plain text password
	// order matters, hashedPassword first, plainTextPassword second
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))

	// Return true if no error (passwords match), false otherwise
	return err == nil
}

// ============================================================================
// JWT TOKEN FUNCTIONS
// ============================================================================

// CreateToken generates a JWT token for an authenticated user.
// This token can be used for subsequent API requests to prove authentication.
//
// Parameters:
//   - user: The authenticated user to create token for
//
// Returns:
//   - string: The JWT token as a string
//
// Token Structure:
// - Algorithm: HS256 (HMAC with SHA-256)
// - Expiration: 1 hour from creation
// - Claims: username and expiration timestamp
//
// Security Features:
// - Short expiration time (1 hour)
// - Signed with secret key
// - Contains user identification
//
// Usage: Called after successful user login
func CreateToken(user User) string {
	// Get current time for token expiration calculation
	now := time.Now()

	// Set token to expire in 1 hour
	// Unix timestamp format for JWT standard
	validUntil := now.Add(time.Hour * 1).Unix()

	// Define the JWT claims (payload)
	// These are the data stored in the token
	claims := jwt.MapClaims{
		"user":    user.Username, // User identification
		"expires": validUntil,    // Expiration timestamp
	}

	// Create the JWT token with HS256 signing method
	// HS256 is HMAC with SHA-256 (symmetric encryption)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	// Use the global JWT secret constant
	secret := BT_JWT_SECRET

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		// Log the error and return empty string
		// In production, you might want to return an error instead
		fmt.Printf("Error signing token: %v\n", err)
		return "N/A"
	}

	return tokenString
}

// ValidateToken validates a JWT token and returns the claims if valid.
// This function is used by middleware to verify authentication tokens.
//
// Parameters:
//   - tokenString: The JWT token to validate
//
// Returns:
//   - jwt.MapClaims: The token claims if valid
//   - error: Any error that occurred during validation
//
// Validation Process:
// - Parses the token string
// - Verifies the signature using the secret key
// - Checks token expiration
// - Returns claims if all validations pass
//
// Usage: Called by middleware for token validation

// REVIEW - AI UPDATE
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token with the same secret used for signing
	// Use the global JWT secret constant
	secret := BT_JWT_SECRET

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	// Check if token is expired
	if exp, ok := claims["expires"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return claims, nil
}

// ============================================================================
// COMPREHENSIVE CODE ANALYSIS
// ============================================================================

/*
EXPLAINING THE CODE:
=========================
This types package implements the core data structures and utility functions
for a secure user authentication system in AWS Lambda. It follows Go best
practices for type safety, security, and error handling.

STRUCTURE OF THE CODE:
===========================
1. Data Structures:
   - RegisterUser: Input validation and serialization
   - User: Internal data representation with security

2. Factory Functions:
   - NewUser: Secure user creation with password hashing

3. Validation Functions:
   - ValidatePassword: Secure password verification

4. JWT Functions:
   - CreateToken: Token generation for authentication

FUNCTION RELATIONSHIPS:
========================
Registration Flow:
RegisterUser → NewUser → User (stored in database)

Login Flow:
User input → ValidatePassword → CreateToken → JWT returned

Authentication Flow:
JWT token → ValidateToken → User access granted

WHY WE USE THIS PATTERN:
========================
• Security: Passwords are never stored in plain text
• Type Safety: Go's type system prevents common errors
• Separation of Concerns: Each function has a single responsibility
• Testability: Functions can be tested independently
• Maintainability: Clear structure makes code easy to understand

HOW IT WORKS:
=============
1. User Registration:
   - Client sends RegisterUser data
   - NewUser hashes the password
   - User struct is stored in database

2. User Login:
   - Client sends username/password
   - ValidatePassword checks credentials
   - CreateToken generates JWT if valid

3. API Access:
   - Client sends JWT in Authorization header
   - Middleware validates JWT
   - Access granted if token is valid

SECURITY CONSIDERATIONS:
========================
• Password Hashing: bcrypt with cost factor 10
• JWT Security: Short expiration time (1 hour)
• Secret Management: TODO - move to environment variables
• Input Validation: TODO - add validation for username/password
• Error Handling: Don't expose internal errors to clients

IMPROVEMENT IDEAS:
==================
• Add input validation for username/password
• Use environment variables for JWT secret
• Add token refresh functionality
• Implement password strength requirements
• Add rate limiting for login attempts
• Use AWS Secrets Manager for sensitive data
• Add comprehensive logging
• Implement token blacklisting for logout

PERFORMANCE OPTIMIZATIONS:
==========================
• Cache bcrypt cost factor in configuration
• Use connection pooling for database
• Implement token caching
• Add request/response compression
• Use CDN for static assets

CONCLUSION:
===========
This types package provides a solid foundation for secure user authentication
in AWS Lambda. The code follows Go best practices and implements proper
security measures for password handling and JWT token generation.

Key Takeaways:
- Always hash passwords before storing
- Use strong, random secret keys for JWT
- Implement proper error handling
- Follow the principle of least privilege
- Test security measures thoroughly
- Keep secrets in environment variables
- Use short token expiration times
- Implement proper input validation
*/
