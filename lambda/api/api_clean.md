/*
===============================================================================
AWS LAMBDA FUNCTION - API HANDLER LAYER
===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This is the API handler layer that processes HTTP requests from API Gateway

	and coordinates with the database layer to implement user authentication
	and registration functionality.

LEARNING OBJECTIVES:
- API Gateway request/response handling
- Business logic implementation
- Error handling and HTTP status codes
- Dependency injection patterns
- Integration with database layer

TOPICS COVERED:
- HTTP request processing
- JSON serialization/deserialization
- User registration workflow
- User authentication workflow
- Error response formatting
- Database integration

FEATURES:
- User registration endpoint (/register)
- User login endpoint (/login)
- Comprehensive error handling
- Input validation
- Database integration
- JWT token generation

ARCHITECTURE:
- ApiHandler struct with database dependency injection
- RegisterUser: Handles new user creation
- LoginUser: Handles user authentication
- Integration with types package for data structures
- Integration with database package for persistence

FUNCTION RELATIONSHIPS & CALL FLOW:
====================================
1. NewApiHandler → ApiHandler (Constructor with dependency injection)
2. RegisterUser → DoesUserExist → NewUser → InsertUser (Registration flow)
3. LoginUser → GetUser → ValidatePassword → CreateToken (Authentication flow)

SECTION ORDER ANALYSIS (REORGANIZED BY LOGICAL FLOW):
 1. Constructor & Dependency Injection (NewApiHandler)
 2. User Registration (RegisterUser)
 3. User Authentication (LoginUser)
 4. Documentation & Analysis

DEVELOPMENT ORDER:
==================
This is the FOURTH function to implement in the development sequence:
1. Database layer (database.go) - FIRST
2. Types layer (types.go) - SECOND
3. Middleware layer (middleware.go) - THIRD
4. API layer (api.go) - FOURTH (THIS FILE)
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST

DYNAMODB V2 MIGRATION NOTES:
============================
This file has been updated to work with AWS SDK v2. Key changes:
- All database calls now require context.Context as first parameter
- Error handling updated for v2 SDK patterns
- TODO comments indicate where context parameters were added
- Database operations use new v2 client patterns

INTEGRATION POINTS:
===================
- Uses database.UserStore interface for database operations
- Uses types.RegisterUser and types.User for data structures
- Returns events.APIGatewayProxyResponse for API Gateway integration
- Handles JSON request/response serialization
*/
package api

import (
	"BTgoAWSlambda/database"
	"BTgoAWSlambda/types"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// !SECTION 1 - CONSTRUCTOR & DEPENDENCY INJECTION
// ApiHandler represents the API handler with database dependency injection.
// This struct follows the dependency injection pattern, allowing for easy
// testing and flexibility in database implementations.
//
// Fields:
//   - dbStore: Interface for database operations (database.UserStore)
//
// Usage: Created by NewApiHandler(), used by main.go for request routing
type ApiHandler struct {
	dbStore database.UserStore
}

// NewApiHandler creates a new ApiHandler instance with dependency injection.
// This constructor follows the dependency injection pattern, allowing the API handler
// to work with any database implementation that implements the UserStore interface.
//
// DEPENDENCY INJECTION EXPLANATION:
// =================================
// Dependency Injection is a design pattern that allows for easy testing and flexibility
// in database implementations. Instead of the API handler creating its own database connection,
// we inject the database dependency from outside. This makes the code more testable and flexible.
//
// BENEFITS:
// =========
// - Testability: Can inject mock database for unit testing
// - Flexibility: Can switch between DynamoDB, PostgreSQL, etc.
// - Separation of Concerns: API layer doesn't manage database connections
// - Interface-based: Uses UserStore interface, not concrete implementation
//
// Parameters:
//   - dbStore: Database interface for user operations (database.UserStore)
//
// Returns:
//   - ApiHandler: Configured handler ready for use
//
// Usage: Called during application initialization in app.go
// Example: apiHandler := NewApiHandler(databaseClient)
func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore, // Inject database dependency
	}
}

// !SECTION 2 - USER REGISTRATION
// RegisterUser handles user registration requests from API Gateway.
// This function processes HTTP POST requests to the /register endpoint,
// validates input data, checks for existing users, and creates new user accounts.
//
// REGISTRATION WORKFLOW:
// =====================
// 1. Deserialize JSON request body to RegisterUser struct
// 2. Validate required fields (username, password)
// 3. Check if user already exists in database
// 4. Create new user with hashed password
// 5. Store user in database
// 6. Return success response
//
// SECURITY FEATURES:
// ==================
// - Input validation for required fields
// - Duplicate user prevention
// - Password hashing (handled by types.NewUser)
// - Error handling without information leakage
//
// Parameters:
//   - request: APIGatewayProxyRequest containing HTTP request data
//
// Returns:
//   - APIGatewayProxyResponse: HTTP response
//   - error: Any error that occurred during processing
//
// HTTP Status Codes:
//   - 200 OK: User successfully registered
//   - 400 Bad Request: Invalid input data
//   - 409 Conflict: User already exists
//   - 500 Internal Server Error: Database or processing error
//
// Usage: Called by main.go when /register endpoint is accessed
func (api *ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Step 1: Deserialize JSON request body to RegisterUser struct
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// Step 2: Validate required fields are not empty
	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Step 3: Check if user already exists
	// DYNAMODB V2 NOTE: context.TODO() is used here - in production,
	// you should pass a proper context with timeout and cancellation
	userExists, err := api.dbStore.DoesUserExist(context.TODO(), registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	// Step 4: Create new user with hashed password
	// This calls types.NewUser() which handles password hashing
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// Step 5: Store user in database
	// DYNAMODB V2 NOTE: context.TODO() is used here - in production,
	// you should pass a proper context with timeout and cancellation
	err = api.dbStore.InsertUser(context.TODO(), user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// Step 6: Return success response
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("User %s registered successfully", registerUser.Username),
		StatusCode: http.StatusOK,
	}, nil
}

// !SECTION 3 - USER AUTHENTICATION
// LoginUser handles user authentication requests from API Gateway.
// This function processes HTTP POST requests to the /login endpoint,
// validates credentials, and returns a JWT token for authenticated users.
//
// AUTHENTICATION WORKFLOW:
// ========================
// 1. Deserialize JSON request body to LoginRequest struct
// 2. Retrieve user data from database by username
// 3. Validate password using bcrypt comparison
// 4. Generate JWT token for authenticated user
// 5. Return JWT token in response
//
// SECURITY FEATURES:
// ==================
// - Password validation using bcrypt
// - JWT token generation with expiration
// - Error handling without information leakage
// - Secure token-based authentication
//
// Parameters:
//   - request: APIGatewayProxyRequest containing HTTP request data
//
// Returns:
//   - APIGatewayProxyResponse: HTTP response with JWT token
//   - error: Any error that occurred during processing
//
// HTTP Status Codes:
//   - 200 OK: User successfully authenticated, JWT token returned
//   - 400 Bad Request: Invalid input data
//   - 401 Unauthorized: Invalid credentials
//   - 500 Internal Server Error: Database or processing error
//
// Usage: Called by main.go when /login endpoint is accessed
func (api *ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Define local struct for login request data
	// This could be moved to types package for consistency
	type LoginRequest struct {
		Username string `json:"username"` // User's username
		Password string `json:"password"` // User's plain text password
	}

	// Step 1: Deserialize JSON request body to LoginRequest struct
	var loginRequest LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// Step 2: Retrieve user data from database
	// DYNAMODB V2 NOTE: context.TODO() is used here - in production,
	// you should pass a proper context with timeout and cancellation
	user, err := api.dbStore.GetUser(context.TODO(), loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, err
	}

	// Step 3: Validate password using bcrypt
	// This calls types.ValidatePassword() which handles secure password comparison
	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	// Step 4: Generate JWT token for authenticated user
	// This calls types.CreateToken() which generates a secure JWT token
	token := types.CreateToken(user)

	// Step 5: Return JWT token in response
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{\"token\": \"%s\"}", token),
		StatusCode: http.StatusOK,
	}, nil
}

// !SECTION 4 - DOCUMENTATION & ANALYSIS
/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA API HANDLER LAYER
===============================================================================

EXPLAINING THE CODE:
=========================
This api.go file implements the API handler layer for AWS Lambda functions,
providing business logic for user registration and authentication. It demonstrates
modern Go development practices including dependency injection, error handling,
and integration with AWS services.

• API Handler: Business logic layer for HTTP request processing
• Dependency Injection: Clean separation of concerns with database abstraction
• Error Handling: Comprehensive HTTP status codes and error responses
• Security: Password hashing, JWT tokens, and input validation
• AWS Integration: Native API Gateway and Lambda integration

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Standard Go packages (context, encoding/json, fmt, net/http)
   - AWS Lambda events package for API Gateway integration
   - Custom packages (database, types)

2. Data Structures:
   - ApiHandler: Main handler struct with database dependency
   - LoginRequest: Local struct for login request data

3. Handler Functions:
   - NewApiHandler: Constructor with dependency injection
   - RegisterUser: User registration endpoint handler
   - LoginUser: User authentication endpoint handler

4. Integration Points:
   - API Gateway: HTTP request/response handling
   - Database: UserStore interface for data persistence
   - Types: Data structures and utility functions

FUNCTION RELATIONSHIPS & CALL FLOW:
====================================
Registration Flow:
1. Client → API Gateway → Lambda → RegisterUser
2. RegisterUser → DoesUserExist → NewUser → InsertUser
3. InsertUser → Database → Success Response

Authentication Flow:
1. Client → API Gateway → Lambda → LoginUser
2. LoginUser → GetUser → ValidatePassword → CreateToken
3. CreateToken → JWT Generation → Success Response

WHY WE USE THIS PATTERN:
========================
• Separation of Concerns: API layer handles business logic only
• Dependency Injection: Database abstraction for flexibility
• Error Handling: Comprehensive HTTP status codes
• Security: Proper password hashing and JWT tokens
• AWS Integration: Native Lambda and API Gateway support

HOW IT WORKS:
=============
1. Request Processing:
   - API Gateway forwards HTTP requests to Lambda
   - Lambda function routes requests to appropriate handlers
   - Handlers process business logic and database operations
   - Responses are returned to API Gateway

2. Registration Process:
   - Validate input data and check for duplicates
   - Hash password and create user record
   - Store user in database and return success

3. Authentication Process:
   - Validate credentials against database
   - Generate JWT token for authenticated users
   - Return token for subsequent API requests

DYNAMODB V2 MIGRATION BENEFITS:
===============================
This API handler works seamlessly with the updated DynamoDB v2 system:

• Context Awareness: All database operations use context.Context
• Performance: Improved error handling and connection management
• Security: Enhanced data validation and error responses
• Scalability: Better support for high-concurrency scenarios

SECURITY FEATURES:
==================
• Password Hashing: bcrypt for secure password storage
• JWT Tokens: Secure token-based authentication
• Input Validation: Comprehensive request validation
• Error Handling: Secure error messages without information leakage
• Database Security: Context-aware database operations

PERFORMANCE CONSIDERATIONS:
===========================
• Context Management: Proper context usage for database operations
• Error Handling: Efficient error processing and response
• Memory Usage: Optimized data structures and processing
• Database Connections: Connection pooling and reuse

TESTING STRATEGY:
=================
• Unit Tests: Individual handler function testing
• Integration Tests: End-to-end request testing
• Mock Testing: Database interface mocking for unit tests
• Security Tests: Authentication and authorization testing

DEPLOYMENT PROCESS:
===================
1. Code Development: Local development with testing
2. Testing: Unit and integration testing
3. Build: Go compilation for Linux architecture
4. Package: ZIP file creation for Lambda deployment
5. Deploy: CDK deployment to AWS
6. Test: API Gateway and Lambda testing

CONCLUSION:
===========
This API handler demonstrates a production-ready approach to AWS Lambda development
with Go. The dependency injection pattern, comprehensive error handling, and security
features make it suitable for real-world applications.

Key Takeaways:
- Always implement proper dependency injection for testability
- Use comprehensive error handling with appropriate HTTP status codes
- Implement security best practices for authentication
- Plan for context-aware database operations
- Follow Go best practices for maintainable code
- Use interface-based design for flexibility
- Implement proper input validation and error handling
- Plan for security from the beginning of development

DEVELOPMENT ORDER REMINDER:
===========================
This is the FOURTH function to implement:
1. Database layer (database.go) - FIRST ✅
2. Types layer (types.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD ✅
4. API layer (api.go) - FOURTH (THIS FILE) ✅
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST

NEXT STEPS:
===========
1. Implement App layer (app.go) with dependency injection
2. Implement Main entry (main.go) with request routing
3. Test all endpoints with different scenarios
4. Implement comprehensive logging and monitoring
5. Add performance optimization and caching
6. Set up CI/CD pipeline for automated deployment
7. Add comprehensive error handling and recovery
8. Implement security best practices
*/
