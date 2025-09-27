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

// This is the "Handler Layer" - the business logic layer that processes
// HTTP requests from API Gateway and coordinates with the database layer.
// It implements the API endpoints and handles request/response processing.

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
	// Database interface for user operations >>> NOT database.DynamoDBClient >>> see database.go for more details
}

// NewApiHandler creates a new ApiHandler instance with database dependency injection.
// This constructor function implements the dependency injection pattern,
// making the handler testable and flexible.
//
// Parameters:
//   - dbStore: Database interface implementation (usually DynamoDBClient)
//   - since this is the ApiHandler layer, we are passing in the UserStore interface
//   - so that we can easily switch between different database implementations : SQLstore YOUstore etc
//
// Returns:
//   - ApiHandler: Configured handler ready for use
//
// Usage: Called during application initialization in app.go
// Example: apiHandler := NewApiHandler(databaseClient)
func NewApiHandler(dbStore database.UserStore) ApiHandler { // NOT database.DynamoDBClient
	return ApiHandler{
		dbStore: dbStore, // Inject database dependency
	}
}

//REVIEW - DEPENDENCY INJECTION with UserStore interface passed in as parameter for Handler
// Dependency Injection is a design pattern that allows for easy testing and flexibility in database implementations.
// It is a way to inject the database client into the API layer so that we can easily switch between different database implementations
// such as DynamoDB, PostgreSQL, etc.
// >>> in api.go, type ApiHandler struct {dbStore database.UserStore} instead of database.DynamoDBClient
// >>> in app.go, type App struct {ApiHandler api.ApiHandler} instead of ApiHandler api.ApiHandler

/* ✅ >>> TO SWITCH BETWEEN DATABASES: >>>>>> IN api.go & db := database.NewDynamoDB() >>>>>> RIGHT THERE
🔗 Component Relationships: NewApiHandler & dBStore & UserStore & DynamoDB
1. NewApiHandler & dBStore & UserStore >>> passes in the UserStore interface to the NewApiHandler function
Why this design?
Dependency Injection: ApiHandler doesn't create its own database connection
Interface-based: Uses UserStore interface, not concrete DynamoDBClient
Testability: Can inject mock database for testing
Flexibility: Can switch between DynamoDB, PostgreSQL, etc.
2. UserStore Interface
Purpose:
Abstraction: Defines what database operations are needed
Contract: Any database implementation must follow this interface
Decoupling: API layer doesn't know about DynamoDB specifics

✅ Data Flow
1. API Gateway → ApiHandler.RegisterUser()
2. ApiHandler → dbStore.DoesUserExist()
3. DynamoDBClient → DynamoDB.GetItem()
4. DynamoDB → Returns user data
5. ApiHandler → Processes response
6. ApiHandler → Returns HTTP response

📊 Summary:
Component		Type				Purpose
dbStore			Interface			Database operations contract
DynamoDBClient	Implementation		Actual DynamoDB operations
UserStore		Interface			Defines database methods
ApiHandler		Consumer			Uses dbStore for operations
*/

// !SECTION 2 - USER REGISTRATION
// RegisterUser handles user registration requests from API Gateway.
// This function implements the complete user registration workflow including
// input validation, duplicate checking, password hashing, and database storage.
//
// FUNCTION CALL FLOW:
// ===================
// 1. JSON deserialization of request body
// 2. Input validation (username/password not empty)
// 3. Database check for existing user (DoesUserExist)
// 4. Password hashing (types.NewUser)
// 5. Database insertion (InsertUser)
// 6. Success response
//
// Parameters:
//   - request: APIGatewayProxyRequest containing JSON user data
//
// Returns:
//   - APIGatewayProxyResponse: HTTP response with status and body
//   - error: Any error that occurred during processing
//
// HTTP Status Codes:
//   - 200 OK: User successfully registered
//   - 400 Bad Request: Invalid input data
//   - 409 Conflict: User already exists
//   - 500 Internal Server Error: Database or processing error
//
// Usage: Called by main.go when /register endpoint is accessed

// NOTE: the old DynamoDB v1 API handler is in the apiOG.md file = SAME
// but in class we tested with api.dbStore.DoesUserExist(event.Username)
// just to test and then then fix it later with API Gateway in the main.go file
func (api *ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Step 1: Deserialize JSON request body to RegisterUser struct
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	// >>> unmarshal the request body from client to the registerUser struct >>> and also the request body is a string
	// >>> so we need to unmarshal it to the registerUser struct which has username and password fields
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request", //API Gtaeway has body and status code
			StatusCode: http.StatusBadRequest,
		}, err // return the error to the developer
	}

	// Step 2: Validate required fields are not empty
	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",     //API Gtaeway has body and status code
			StatusCode: http.StatusBadRequest, // without API Gateway, we would just return the error
		}, fmt.Errorf("register user request fields cannot be empty")
	}

	// Step 3: Check if user already exists in database
	// DYNAMODB V2 NOTE: context.TODO() is used here - in production,
	// you should pass a proper context with timeout and cancellation
	doesUserExist, err := api.dbStore.DoesUserExist(context.TODO(), registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("checking if user exists error %w", err)
	}

	// Step 4: Return conflict if user already exists
	if doesUserExist {
		return events.APIGatewayProxyResponse{
			Body:       "Hey there. See again. User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	// Step 5: Create user with hashed password
	// This calls types.NewUser which handles bcrypt password hashing
	//NOTE: check password hash in types.go : RegisterUser -> NewUser -> NewUser(registerUser) > passwordHash > user.PasswordHash
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error hashing user password %w", err)
	}

	// FIXME: ADD DEBUG LOGGING HERE
	fmt.Printf("DEBUG: Retrieved user from DB: %+v\n", user)
	fmt.Printf("DEBUG: Stored password hash length: %d\n", len(user.PasswordHash))
	if len(user.PasswordHash) > 0 {
		hashPreview := user.PasswordHash
		if len(hashPreview) > 10 {
			hashPreview = hashPreview[:10]
		}
		fmt.Printf("DEBUG: Stored password hash starts with: %s\n", hashPreview)
	} else {
		fmt.Printf("DEBUG: Password hash is EMPTY!\n")
	}
	fmt.Printf("DEBUG: Login password: %s\n", registerUser.Password)

	// Step 6: Insert user into database
	err = api.dbStore.InsertUser(context.TODO(), user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user into the database %w", err)
	}

	// ADD MORE DEBUG LOGGING
	fmt.Printf("DEBUG: User inserted successfully\n")
	fmt.Printf("DEBUG: Final user object: %+v\n", user)
	fmt.Printf("DEBUG: Final password hash length: %d\n", len(user.PasswordHash))

	// Step 7: Return success response
	return events.APIGatewayProxyResponse{
		Body:       "BINGO !!! Successfully, Registered Voter !",
		StatusCode: http.StatusOK,
	}, nil // return the success response to the client
}

// !SECTION 3 - USER AUTHENTICATION
// LoginUser handles user authentication requests from API Gateway.
// This function implements the complete user login workflow including
// credential validation, password verification, and JWT token generation.
//
// FUNCTION CALL FLOW:
// ===================
// 1. JSON deserialization of request body
// 2. Database retrieval of user data (GetUser)
// 3. Password validation (ValidatePassword)
// 4. JWT token generation (CreateToken)
// 5. Success response with token
//
// Parameters:
//   - request: APIGatewayProxyRequest containing JSON login credentials
//
// Returns:
//   - APIGatewayProxyResponse: HTTP response with JWT token or error
//   - error: Any error that occurred during processing
//
// HTTP Status Codes:
//   - 200 OK: Login successful, returns JWT token
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
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// Step 2: Retrieve user data from database
	// DYNAMODB V2 NOTE: context.TODO() is used here - in production,
	// you should pass a proper context with timeout and cancellation
	//REVIEW - We created the user in the database in the RegisterUser function
	// Now we are retrieving the GetUser from the database by username
	// so that we can validate the password with the password hash in the database
	//  by user.PasswordHash in the types.go file
	// RegisterUser -> NewUser -> NewUser(registerUser) > passwordHash > user.PasswordHash
	// and then generate a JWT token for the authenticated user
	user, err := api.dbStore.GetUser(context.TODO(), loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	//FIXME - Loging Errors - not working Debugging needed
	// ADD DEBUG LOGGING HERE
	fmt.Printf("DEBUG: Retrieved user from DB: %+v\n", user)
	fmt.Printf("DEBUG: Username: %s\n", user.Username)
	fmt.Printf("DEBUG: PasswordHash field: %s\n", user.PasswordHash)
	fmt.Printf("DEBUG: Stored password hash length: %d\n", len(user.PasswordHash))
	if len(user.PasswordHash) > 0 {
		hashPreview := user.PasswordHash
		if len(hashPreview) > 10 {
			hashPreview = hashPreview[:10]
		}
		fmt.Printf("DEBUG: Stored password hash starts with: %s\n", hashPreview)
	} else {
		fmt.Printf("DEBUG: Password hash is EMPTY!\n")
	}
	fmt.Printf("DEBUG: Login password: %s\n", loginRequest.Password)
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Step 3: Validate password using bcrypt comparison
	// This calls types.ValidatePassword which uses bcrypt.CompareHashAndPassword
	// Order of parameters is important: user.PasswordHash, loginRequest.Password
	isValid := types.ValidatePassword(user.PasswordHash, loginRequest.Password)
	fmt.Printf("DEBUG: Password validation result: %v\n", isValid)

	if !isValid {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid login credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	//TODO: - JWT token generation is not implemented yet >> Middleware layer
	// Step 4: Generate JWT token for authenticated user
	// This calls types.CreateToken which creates a signed JWT with user info
	accessToken := types.CreateToken(user)

	// Step 5: Format success response with JWT token
	// Returns JSON with access_token field for client consumption
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body: successMsg, // TODO: Implement JWT token generation
		// Body:       "Almmost successfull Login, JWT token generation is not implemented yet",
		StatusCode: http.StatusOK,
	}, nil
}

/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA API HANDLER LAYER
===============================================================================

EXPLAINING THE CODE:
=========================
This API handler layer implements the business logic for user authentication
and registration in a serverless AWS Lambda function. It follows clean architecture
principles with proper separation of concerns and dependency injection.

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Standard Go packages (context, encoding/json, fmt, net/http)
   - AWS Lambda events package for API Gateway integration
   - Custom packages (database, types) for business logic

2. Data Structures:
   - ApiHandler: Main handler struct with database dependency
   - LoginRequest: Local struct for login request data

3. Handler Functions:
   - NewApiHandler: Constructor with dependency injection
   - RegisterUser: Complete user registration workflow
   - LoginUser: Complete user authentication workflow

4. Integration Points:
   - Database layer: Uses UserStore interface for persistence
   - Types layer: Uses RegisterUser, User, and utility functions
   - API Gateway: Returns proper HTTP responses

FUNCTION RELATIONSHIPS & DEPENDENCIES:
======================================
1. NewApiHandler → ApiHandler (Constructor pattern)
2. RegisterUser → DoesUserExist → NewUser → InsertUser (Registration flow)
3. LoginUser → GetUser → ValidatePassword → CreateToken (Authentication flow)

Layer Dependencies:
- API Layer (THIS) → Database Layer (database.go)
- API Layer (THIS) → Types Layer (types.go)
- Main Layer (main.go) → API Layer (THIS)

WHY WE USE THIS PATTERN:
========================
• Separation of Concerns: Each layer has a specific responsibility
• Dependency Injection: Makes testing easier and code more flexible
• Interface-Based Design: Allows for easy mocking and testing
• Error Handling: Comprehensive error handling with proper HTTP status codes
• Clean Architecture: Business logic separated from infrastructure concerns

HOW IT WORKS:
=============
1. API Gateway receives HTTP request
2. Lambda runtime invokes main() function
3. main() routes request to appropriate handler (RegisterUser/LoginUser)
4. Handler validates input and calls database layer
5. Database layer performs operations with DynamoDB v2
6. Handler processes results and returns HTTP response
7. API Gateway transforms response back to HTTP

DYNAMODB V2 MIGRATION CHANGES:
==============================
This code has been updated for AWS SDK v2 with the following changes:

1. Context Support:
   - All database calls now require context.Context as first parameter
   - Uses context.TODO() (should be replaced with proper context in production)

2. Error Handling:
   - Updated error wrapping patterns for v2 SDK
   - Proper error propagation through call chain

3. Client Patterns:
   - Uses new v2 client initialization patterns
   - Updated method signatures for v2 compatibility

4. Type Safety:
   - Improved type safety with v2 SDK patterns
   - Better error handling and validation

APPLICATION IN USE:
==================
• User Registration: POST /register with {"username": "user", "password": "pass"}
• User Login: POST /login with {"username": "user", "password": "pass"}
• Response Format: JSON with appropriate HTTP status codes
• JWT Token: Returned on successful login for subsequent requests

IMPROVEMENT IDEAS:
==================
• Add input validation middleware
• Implement proper context with timeouts
• Add request/response logging
• Implement rate limiting
• Add comprehensive error handling
• Use structured logging (logrus, zap)
• Add input sanitization
• Implement request timeout handling
• Add health check endpoint
• Use environment variables for configuration

SECURITY CONSIDERATIONS:
========================
• Password Hashing: Handled by types.NewUser with bcrypt
• JWT Security: Tokens generated with proper expiration
• Input Validation: Basic validation implemented
• Error Messages: Don't expose internal details
• HTTPS Enforcement: Should be handled at API Gateway level
• Rate Limiting: Should be implemented to prevent abuse

PERFORMANCE OPTIMIZATIONS:
==========================
• Connection Pooling: Handled by DynamoDB v2 SDK
• Caching: Could add caching layer for frequently accessed data
• Memory Usage: Optimize Lambda memory allocation
• Cold Start: Minimize cold start time with proper initialization
• Database Queries: Optimize DynamoDB operations

CONCLUSION:
===========
This API handler layer provides a solid foundation for user authentication
in AWS Lambda. The code follows Go best practices and implements proper
error handling, dependency injection, and clean architecture principles.

Key Takeaways:
- Always implement proper error handling with appropriate HTTP status codes
- Use dependency injection for testability and flexibility
- Follow the repository pattern for database operations
- Implement proper input validation and sanitization
- Use context for cancellation and timeout handling
- Plan for monitoring and logging from the start
- Consider security implications of every operation
- Test thoroughly before deployment
- Use AWS SDK v2 best practices
- Implement proper separation of concerns

DEVELOPMENT ORDER REMINDER:
===========================
This is the FOURTH function to implement:
1. Database layer (database.go) - FIRST
2. Types layer (types.go) - SECOND
3. Middleware layer (middleware.go) - THIRD
4. API layer (api.go) - FOURTH (THIS FILE)
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST
*/
