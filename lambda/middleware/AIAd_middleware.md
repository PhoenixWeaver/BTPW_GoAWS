/*
===============================================================================
AWS LAMBDA FUNCTION - MIDDLEWARE PACKAGE
===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This package contains middleware functions for the AWS Lambda function,
             including JWT authentication middleware for protected routes.

LEARNING OBJECTIVES:
- Middleware pattern implementation
- JWT token validation
- Cross-cutting concerns handling
- Request/response processing

TOPICS COVERED:
- JWT token parsing and validation
- Middleware function composition
- Error handling in middleware
- HTTP response formatting

FEATURES:
- JWT authentication middleware
- Request validation
- Error response handling
- Token extraction from headers

ARCHITECTURE:
- ValidateJWTMiddleware: JWT token validation middleware
- Middleware function composition pattern
- Integration with types package for JWT operations
*/

package middleware

import (
	"BTgoAWSlambda/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ValidateJWTMiddleware is a middleware function that validates JWT tokens
// for protected routes. It extracts the token from the Authorization header,
// validates it, and either allows the request to proceed or returns an error.
//
// Parameters:
//   - handler: The next handler function to call if validation succeeds
//
// Returns:
//   - A new handler function that includes JWT validation
//
// Usage: Wrap protected handlers with this middleware
// Example: middleware.ValidateJWTMiddleware(ProtectedHandler)

func ValidateJWTMiddleware(handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Extract Authorization header
		authHeader := request.Headers["Authorization"]
		if authHeader == "" {
			authHeader = request.Headers["authorization"] // Try lowercase
		}

		// Check if Authorization header exists
		if authHeader == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Authorization header required",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		// Extract token from "Bearer <token>" format
		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		// Validate the JWT token
		_, err := types.ValidateToken(token)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "Invalid token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		// Token is valid, proceed to the actual handler
		return handler(request)
	}
}

/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA MIDDLEWARE
===============================================================================

EXPLAINING THE CODE:
=========================
This middleware package implements the middleware pattern for AWS Lambda functions,
providing JWT authentication for protected routes.

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Standard Go packages (net/http)
   - AWS Lambda events package
   - Custom types package for JWT operations

2. Middleware Functions:
   - ValidateJWTMiddleware: JWT token validation middleware
   - Function composition pattern for middleware chaining

3. Integration Points:
   - Uses types.ValidateToken for JWT validation
   - Returns proper HTTP responses for API Gateway
   - Handles Authorization header extraction

WHY WE USE THIS PATTERN:
========================
• Separation of Concerns: Authentication logic separated from business logic
• Reusability: Middleware can be applied to multiple handlers
• Function Composition: Clean way to chain middleware functions
• Error Handling: Consistent error responses for authentication failures

HOW IT WORKS:
=============
1. Middleware function wraps the original handler
2. Extracts Authorization header from request
3. Validates JWT token using types.ValidateToken
4. Either calls original handler or returns error response
5. Maintains proper HTTP status codes for API Gateway

APPLICATION IN USE:
==================
• Protected Routes: Applied to handlers that require authentication
• JWT Validation: Ensures only valid tokens can access protected resources
• Error Responses: Returns 401 Unauthorized for invalid/missing tokens

IMPROVEMENT IDEAS:
==================
• Add request logging middleware
• Implement rate limiting middleware
• Add CORS middleware
• Implement request validation middleware
• Add metrics and monitoring middleware

SECURITY CONSIDERATIONS:
========================
• JWT Token Validation: Proper token verification
• Header Extraction: Secure token extraction from headers
• Error Messages: Don't expose sensitive information
• Token Format: Support for Bearer token format

PERFORMANCE OPTIMIZATIONS:
==========================
• Token Caching: Cache validated tokens for performance
• Early Return: Return errors quickly for invalid requests
• Minimal Processing: Keep middleware lightweight

CONCLUSION:
===========
This middleware package provides a clean, reusable way to implement
authentication in AWS Lambda functions using the middleware pattern.

Key Takeaways:
- Use middleware for cross-cutting concerns
- Implement proper error handling and HTTP status codes
- Keep middleware functions lightweight and focused
- Use function composition for middleware chaining
- Validate tokens securely and consistently
*/
