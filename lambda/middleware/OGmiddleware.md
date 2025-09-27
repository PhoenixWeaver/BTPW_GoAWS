/*
===============================================================================
AWS LAMBDA FUNCTION - MIDDLEWARE PACKAGE (JWT AUTHENTICATION)
===============================================================================

DEVELOPMENT ORDER: THIRD FUNCTION TO IMPLEMENT
===============================================
This is the THIRD function in the development sequence:
1. Types layer (types.go) - FIRST ✅
2. Database layer (database.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD (THIS FILE) ✅
4. API layer (api.go) - FOURTH
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST

Author: Ben Tran
Date: 22/09/2025
Description: This package implements JWT authentication middleware for AWS Lambda
             functions. It provides secure token validation and user authentication
             for protected routes in serverless applications.

LEARNING OBJECTIVES:
- Middleware pattern implementation in Go
- JWT token validation and parsing
- AWS Lambda middleware integration
- Cross-cutting concerns (authentication)
- Error handling in middleware

TOPICS COVERED:
- JWT token extraction from HTTP headers
- Token parsing and validation
- Middleware pattern for authentication
- Error handling and response formatting
- AWS Lambda event processing

FEATURES:
- JWT token validation middleware
- Authorization header parsing
- Token expiration checking
- Secure error handling
- AWS Lambda integration

ARCHITECTURE:
- ValidateJWTMiddleware: Main middleware function
- extractTokenFromHeaders: Token extraction utility
- parseToken: JWT parsing and validation
- Integration with golang-jwt/jwt/v5 library

FUNCTION RELATIONSHIPS & CALL FLOW:
====================================
1. ValidateJWTMiddleware → extractTokenFromHeaders → parseToken (Authentication flow)
2. main.go → ValidateJWTMiddleware → ProtectedHandler (Protected route flow)
3. API Gateway → Lambda → Middleware → Handler (Request processing flow)

DYNAMODB V2 MIGRATION NOTES:
============================
This middleware package works with the updated DynamoDB v2 system:
- JWT tokens are validated before database operations
- Middleware runs before API handlers
- No direct database dependencies in middleware
- Works with context-aware database operations

INTEGRATION POINTS:
===================
- Uses golang-jwt/jwt/v5 for JWT operations
- Integrates with AWS Lambda events package
- Works with API Gateway proxy requests
- Returns proper HTTP responses for authentication failures

 Section Order Analysis:
  1. ValidateJWTMiddleware (main function): Main middleware function with extensive documentation
   ↓ calls
  2. extractTokenFromHeaders (utility): Token extraction utility function
   ↓ calls
  3. parseToken (utility): JWT parsing and validation utility function
*/

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"BTgoAWSlambda/types" // Import the types package for JWT secret constant and CreateToken()

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// !SECTION 1 - ValidateJWTMiddleware
// ValidateJWTMiddleware is a higher-order function that implements JWT authentication middleware.
// This function follows the middleware pattern, wrapping handlers with authentication logic.
//
// MIDDLEWARE PATTERN EXPLANATION:
// ===============================
// Middleware is a design pattern that allows you to add cross-cutting concerns (like authentication)
// to your application without modifying the core business logic. It acts as a "wrapper" around handlers.
//
// FUNCTION SIGNATURE:
// ===================
// Input:  next function (the actual handler to wrap)
// Output: wrapped function (handler with authentication)
//
// HOW IT WORKS:
// =============
// 1. Receives a request from API Gateway
// 2. Extracts JWT token from Authorization header
// 3. Validates the token (signature, expiration, format)
// 4. If valid: calls the next handler
// 5. If invalid: returns 401 Unauthorized response
//
// USAGE EXAMPLE:
// ==============
// protectedHandler := middleware.ValidateJWTMiddleware(ProtectedHandler)
// response := protectedHandler(request)
//
// DYNAMODB V2 INTEGRATION:
// ========================
// This middleware works with the updated DynamoDB v2 system by validating JWT tokens
// before any database operations occur. It ensures only authenticated users can access
// protected resources that require database access.
//
// Parameters:
//   - next: The handler function to wrap with authentication
//
// Returns:
//   - func: A new handler function that includes JWT validation
//
// HTTP Status Codes:
//   - 401 Unauthorized: Missing token, invalid token, or expired token
//   - 200 OK: Token is valid, proceeds to next handler
//
// Security Features:
//   - Token extraction from Authorization header
//   - JWT signature validation
//   - Token expiration checking
//   - Proper error handling without exposing internal details
func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// REVIEW: Why do we return exact same func() as the next handler?
		// Because we want to wrap the next handler with authentication logic
		// and return the same type of function as the next handler
		// so that we can chain middleware functions
		// and return the same type of function as the next handler
		// and so on

		// STEP 1: Extract JWT token from Authorization header
		// This calls extractTokenFromHeaders() which parses the "Authorization: Bearer <token>" header
		// and returns just the token part, or empty string if not found
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			// No token found in request headers
			return events.APIGatewayProxyResponse{
				Body:       "Missing Auth token",    // User-friendly error message
				StatusCode: http.StatusUnauthorized, // HTTP 401 Unauthorized
			}, nil
		}

		// STEP 2: Parse and validate the JWT token
		// This calls parseToken() which:
		// - Parses the JWT token string
		// - Validates the signature using the secret key
		// - Extracts the claims (user data) from the token
		// - Returns error if token is malformed or invalid
		claims, err := parseToken(tokenString)
		if err != nil {
			// Token is invalid (malformed, wrong signature, etc.)
			return events.APIGatewayProxyResponse{
				Body:       "User Unauthorized",     // User-friendly error message
				StatusCode: http.StatusUnauthorized, // HTTP 401 Unauthorized
			}, err // Return the error for logging purposes
		}

		// STEP 3: Check if the token has expired
		// JWT tokens contain an "expires" claim with a Unix timestamp
		// We compare this with the current time to check if the token is still valid
		expires := int64(claims["expires"].(float64)) // Extract expiration timestamp from claims
		if time.Now().Unix() > expires {
			// Token has expired
			return events.APIGatewayProxyResponse{
				Body:       "token expired",         // User-friendly error message
				StatusCode: http.StatusUnauthorized, // HTTP 401 Unauthorized
			}, nil
		}

		// STEP 4: Token is valid - proceed to the next handler
		// At this point, we know:
		// - Token exists and is properly formatted
		// - Token signature is valid
		// - Token has not expired
		// - User is authenticated and authorized
		return next(request) // Call the actual handler (e.g., ProtectedHandler)
	}
}

// !SECTION 2 - extractTokenFromHeaders
// extractTokenFromHeaders extracts the JWT token from the Authorization header.
// This function parses the standard "Authorization: Bearer <token>" header format
// and returns just the token part for JWT validation.
//
// AUTHORIZATION HEADER FORMAT:
// ============================
// Standard format: "Authorization: Bearer <jwt-token>"
// Example: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//
// HOW IT WORKS:
// =============
// 1. Looks for "Authorization" key in request headers
// 2. Splits the header value by "Bearer " prefix
// 3. Returns the token part (after "Bearer ")
// 4. Returns empty string if header is missing or malformed
//
// SECURITY CONSIDERATIONS:
// ========================
// - Validates header format to prevent injection attacks
// - Returns empty string for malformed headers
// - Case-sensitive header key matching
// - Proper Bearer token format validation
//
// Parameters:
//   - headers: Map of HTTP headers from API Gateway request
//
// Returns:
//   - string: JWT token if found and properly formatted, empty string otherwise
//
// Usage: Called by ValidateJWTMiddleware to extract token from request
//
// DYNAMODB V2 INTEGRATION:
// ========================
// This function works with the updated DynamoDB v2 system by providing clean
// token extraction that works with the context-aware authentication flow.
func extractTokenFromHeaders(headers map[string]string) string {
	// STEP 1: Check if Authorization header exists
	// API Gateway passes headers as a map[string]string
	// We need to check if the "Authorization" key exists
	authHeader, ok := headers["Authorization"]

	if !ok {
		// No Authorization header found in request
		return "" // Return empty string to indicate no token
	}

	// STEP 2: Parse the Bearer token format
	// Standard format: "Authorization: Bearer <token>"
	// We split by "Bearer " to separate the prefix from the actual token
	splitToken := strings.Split(authHeader, "Bearer ") //NOTE - Make sure there is a space after "Bearer "
	if len(splitToken) != 2 {
		// Header doesn't follow "Bearer <token>" format
		// This could be:
		// - "Authorization: <token>" (missing Bearer prefix)
		// - "Authorization: Bearer" (missing token)
		// - "Authorization: Bearer <token> <extra>" (malformed)
		return "" // Return empty string for malformed header
	}

	// STEP 3: Return the actual JWT token
	// splitToken[0] = "" (empty string before "Bearer ")
	// splitToken[1] = "<actual-jwt-token>" (the token we want)
	return splitToken[1] // Return the JWT token for validation
}

// !SECTION 3 - parseToken - JWT PARSING AND VALIDATION
// parseToken parses and validates a JWT token string.
// This function performs comprehensive JWT validation including signature verification,
// format validation, and claims extraction for authentication purposes.
//
// JWT VALIDATION PROCESS:
// ======================
// 1. Parse the JWT token string into its components (header, payload, signature)
// 2. Verify the signature using the secret key
// 3. Check token format and structure
// 4. Extract and return the claims (user data)
//
// SECURITY FEATURES:
// ==================
// - Signature verification using HMAC-SHA256
// - Token format validation
// - Claims extraction and validation
// - Proper error handling without exposing internal details
//
// JWT STRUCTURE:
// ==============
// JWT tokens have three parts separated by dots:
// 1. Header: Algorithm and token type
// 2. Payload: Claims (user data, expiration, etc.)
// 3. Signature: HMAC signature for verification
//
// Parameters:
//   - tokenString: The JWT token string to parse and validate
//
// Returns:
//   - jwt.MapClaims: The token claims (user data) if valid
//   - error: Any error that occurred during parsing or validation
//
// Usage: Called by ValidateJWTMiddleware to validate JWT tokens
//
// DYNAMODB V2 INTEGRATION:
// ========================
// This function works with the updated DynamoDB v2 system by providing secure
// JWT validation that works with the context-aware authentication flow.
func parseToken(tokenString string) (jwt.MapClaims, error) {
	// STEP 1: Parse the JWT token with signature verification
	// jwt.Parse() performs the following validations:
	// - Token format validation (three parts separated by dots)
	// - Header parsing and algorithm verification
	// - Signature verification using the secret key
	// - Claims extraction and validation
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // interface{} is a placeholder for any type of data so we can return the secret key
		// Use the global JWT secret constant from types package
		secret := types.BT_JWT_SECRET
		return []byte(secret), nil // Return the secret key for signature verification
	})

	// STEP 2: Handle parsing errors
	if err != nil {
		// Token parsing failed (malformed, wrong signature, etc.)
		// Return generic error message to avoid exposing internal details
		return nil, fmt.Errorf("unauthorized") // Generic error for security
	}

	// STEP 3: Verify token validity
	if !token.Valid {
		// Token is not valid (expired, malformed, etc.)
		return nil, fmt.Errorf("token is not valid - unauthorized")
	}

	// STEP 4: Extract claims from the token
	// Claims contain the user data (username, expiration, etc.)
	claims, ok := token.Claims.(jwt.MapClaims) // from types.go - CreateToken() function we can see that the claims are of type jwt.MapClaims so we need to cast the token.Claims to jwt.MapClaims
	if !ok {
		// Claims extraction failed (unexpected token type)
		return nil, fmt.Errorf("token of unrecognized type - unauthorized")
	}

	// STEP 5: Return the validated claims
	// At this point, we know:
	// - Token is properly formatted
	// - Signature is valid
	// - Token is not expired
	// - Claims are accessible
	return claims, nil // Return the claims for further processing
}

/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA MIDDLEWARE PACKAGE
===============================================================================

EXPLAINING THE CODE:
=========================
This middleware package implements JWT authentication for AWS Lambda functions
using the middleware pattern. It provides secure token validation and user
authentication for protected routes in serverless applications.

• Middleware Pattern: Wraps handlers with authentication logic
• JWT Validation: Comprehensive token parsing and verification
• Security Features: Signature verification, expiration checking, error handling
• AWS Integration: Works with API Gateway and Lambda events

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Standard Go packages (fmt, net/http, strings, time)
   - AWS Lambda events package for API Gateway integration
   - JWT library (golang-jwt/jwt/v5) for token operations

2. Middleware Functions:
   - ValidateJWTMiddleware: Main middleware wrapper function
   - extractTokenFromHeaders: Token extraction utility
   - parseToken: JWT parsing and validation

3. Integration Points:
   - API Gateway: Processes HTTP requests with headers
   - Lambda Events: Uses APIGatewayProxyRequest/Response types
   - JWT Library: Handles token parsing and validation

FUNCTION RELATIONSHIPS & CALL FLOW:
====================================
Authentication Flow:
1. API Gateway → Lambda → ValidateJWTMiddleware
2. ValidateJWTMiddleware → extractTokenFromHeaders
3. extractTokenFromHeaders → parseToken
4. parseToken → JWT validation
5. Valid token → Next handler
6. Invalid token → 401 Unauthorized

Protected Route Flow:
1. Client → API Gateway → Lambda
2. Lambda → ValidateJWTMiddleware(ProtectedHandler)
3. Middleware → Token validation
4. Valid → ProtectedHandler → Response
5. Invalid → 401 Unauthorized

WHY WE USE THIS PATTERN:
========================
• Middleware Pattern: Separates authentication from business logic
• JWT Security: Industry-standard token-based authentication
• Reusability: Can be applied to any handler function
• Security: Comprehensive token validation and error handling
• AWS Integration: Works seamlessly with Lambda and API Gateway

HOW IT WORKS:
=============
1. Request Processing:
   - API Gateway receives HTTP request with Authorization header
   - Lambda function invokes ValidateJWTMiddleware
   - Middleware extracts JWT token from headers
   - Token is parsed and validated for signature and expiration
   - Valid token allows access to protected handler
   - Invalid token returns 401 Unauthorized response

2. JWT Validation Process:
   - Token format validation (three parts separated by dots)
   - Signature verification using secret key
   - Claims extraction and validation
   - Expiration timestamp checking
   - Error handling for malformed or invalid tokens

3. Security Features:
   - Bearer token format validation
   - HMAC-SHA256 signature verification
   - Token expiration checking
   - Proper error handling without information leakage

DYNAMODB V2 MIGRATION BENEFITS:
===============================
This middleware package works seamlessly with the updated DynamoDB v2 system:

• Context Awareness: Middleware validates tokens before database operations
• Performance: JWT validation happens before expensive database calls
• Security: Ensures only authenticated users access database resources
• Integration: Works with context-aware database operations

APPLICATION IN USE:
==================
• Protected Routes: /protected endpoint requires valid JWT token
• Authentication: Validates user identity before allowing access
• Authorization: Ensures only authenticated users can access resources
• Error Handling: Returns appropriate HTTP status codes for different scenarios

IMPROVEMENT IDEAS:
==================
• Add token refresh mechanism for long-lived sessions
• Implement role-based access control (RBAC)
• Add request rate limiting to prevent abuse
• Use AWS Secrets Manager for JWT secret keys
• Implement token blacklisting for logout functionality
• Add comprehensive logging for security monitoring
• Use environment variables for configuration
• Implement token rotation for enhanced security
• Add support for multiple JWT algorithms
• Implement token caching for performance

SECURITY CONSIDERATIONS:
========================
• JWT Secret Management: Move secrets to environment variables
• Token Expiration: Implement proper expiration handling
• Error Handling: Don't expose internal implementation details
• Input Validation: Validate all token components
• Rate Limiting: Prevent brute force attacks
• Logging: Monitor authentication attempts and failures
• HTTPS Enforcement: Ensure all communications are encrypted
• Token Storage: Consider secure token storage on client side

PERFORMANCE OPTIMIZATIONS:
==========================
• Token Caching: Cache validated tokens to reduce parsing overhead
• Connection Pooling: Optimize database connections
• Memory Usage: Minimize memory allocation in token parsing
• Cold Start: Optimize Lambda cold start time
• Request Batching: Process multiple requests efficiently

TESTING STRATEGY:
=================
• Unit Tests: Test each middleware function independently
• Integration Tests: Test complete authentication flow
• Security Tests: Test token validation and error handling
• Load Tests: Test middleware performance under load
• Mock Tests: Test with mock JWT tokens and invalid inputs

CONCLUSION:
===========
This middleware package provides a solid foundation for JWT authentication
in AWS Lambda functions. The middleware pattern implementation, comprehensive
security features, and AWS integration make it suitable for production use.

Key Takeaways:
- Always implement proper JWT validation with signature verification
- Use middleware pattern for cross-cutting concerns like authentication
- Implement comprehensive error handling without exposing internal details
- Plan for security from the beginning of development
- Test authentication thoroughly before deployment
- Use environment variables for sensitive configuration
- Implement proper logging and monitoring for security
- Consider performance implications of authentication
- Follow the exact development order: Types → Database → Middleware → API → App → Main
- Use AWS SDK v2 best practices for better performance and features

DEVELOPMENT ORDER REMINDER:
===========================
This is the THIRD function to implement:
1. Types layer (types.go) - FIRST ✅
2. Database layer (database.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD (THIS FILE) ✅
4. API layer (api.go) - FOURTH
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST

NEXT STEPS:
===========
1. Implement API layer (api.go) with business logic
2. Implement App layer (app.go) with dependency injection
3. Implement Main entry (main.go) with request routing
4. Set up CDK infrastructure (cdk7.go)
5. Add comprehensive testing
6. Deploy and test the complete application
*/
