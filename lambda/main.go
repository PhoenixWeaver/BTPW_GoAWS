/*
===============================================================================
AWS LAMBDA FUNCTION - GO LEARNING PROJECT
===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This is a comprehensive AWS Lambda function built with Go that demonstrates
             user authentication, middleware patterns, and API Gateway integration.


LEARNING OBJECTIVES:
- AWS Lambda function development with Go
- API Gateway integration and request handling
- Middleware pattern implementation
- JWT authentication and protected routes
- Error handling and response formatting
- Dependency injection and application structure

TOPICS COVERED:
- AWS Lambda Go runtime
- API Gateway proxy integration
- Middleware pattern for cross-cutting concerns
- User registration and authentication
- Protected route implementation
- HTTP status codes and response handling

FEATURES:
- User registration endpoint (/register)
- User login endpoint (/login)
- Protected route with JWT middleware (/protected)
- Comprehensive error handling
- Modular application structure
- Middleware-based authentication

ARCHITECTURE:
- main.go: Entry point and request routing
- app/: Application structure and dependency injection
- middleware/: Cross-cutting concerns (auth, logging)
- handlers/: Business logic and request processing
- models/: Data structures and contracts
- utils/: Helper functions and response formatting

SECTION ORDER ANALYSIS (REORGANIZED BY LOGICAL FLOW):
  1. Data Structures & Types (MyEvent struct)
  2. Utility Functions (HandleRequest, ProtectedHandler)
  3. Main Entry Point (main function with routing logic)
  4. Documentation & Analysis


/// !SECTION 1 - DATA STRUCTURES & TYPES
/// !SECTION 2 - UTILITY FUNCTIONS
/// !SECTION 3 - MAIN ENTRY POINT
/// !SECTION 4 - DOCUMENTATION & ANALYSIS
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"BTgoAWSlambda/app"
	"BTgoAWSlambda/database"
	"BTgoAWSlambda/middleware"
	"BTgoAWSlambda/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// =============================================================
// !SECTION 1 - DATA STRUCTURES & TYPES
// =============================================================
// MyEvent represents the input data structure for our Lambda function.
// This struct defines the expected format for incoming requests when testing
// the Lambda function directly (not through API Gateway).
//
// STRUCTURE:
// ==========
// - Username: User's chosen username (required)
// - Password: User's plain text password (will be hashed)
//
// USAGE:
// ======
// This struct is used for direct Lambda function testing and local development.
// When using API Gateway, the request data comes through the APIGatewayProxyRequest
// structure instead.
//
// TESTING:
// ========
// You can test this function by sending JSON data like:
// {"username": "testuser", "password": "testpass"}
type MyEvent struct {
	Username string `json:"username"` // User's username (required)
	Password string `json:"password"` // User's password (required)
}

// =============================================================
// !SECTION 2 - UTILITY FUNCTIONS
// =============================================================
// HandleRequest is a utility function for direct Lambda function testing.
// This function demonstrates basic Lambda function structure and error handling.
//
// Parameters:
//   - event: MyEvent struct containing the input data
//
// Returns:
//   - string: Success message with username
//   - error: Error if validation fails
//
// Usage: This function can be used for direct Lambda invocation testing
func HandleRequest(event MyEvent) (string, error) {
	// Validate that username is not empty
	if event.Username == "" {
		return "Invalid Request", fmt.Errorf("Yo, Anonymous ? Username cannot be empty !")
	}
	// Validate that password is not empty
	if event.Password == "" {
		return "Invalid Request", fmt.Errorf("Yo, Emo ? Password cannot be empty !")
	}

	//TODO -DATABASE OPERATIONS
	// Create database connection with appropriate table name
	var tableName string
	// Check if running locally (no AWS_LAMBDA_FUNCTION_NAME environment variable)
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		// if AWS_LAMBDA_FUNCTION_NAME is not set, then we are running locally with "LocoTable"
		tableName = "LocoTable" // Local testing table
	} else {
		// elseif There is Lambda Function, then we are running on AWS with "BTtableGuestsInfo"
		// (either live or test environment - same table name): test Lambda or live AWS
		tableName = os.Getenv("TABLE_NAME") // AWS table: "BTtableGuestsInfo"
		if tableName == "" {
			tableName = "BTtableGuestsInfo" // Default AWS table
		}
	}

	dbClient := database.NewDynamoDBWithTable(tableName)

	// ✅ FIX: Create RegisterUser struct first
	registerUser := types.RegisterUser{
		Username: event.Username,
		Password: event.Password,
	}

	// ✅ FIX: Use NewUser to hash the password
	user, err := types.NewUser(registerUser)
	if err != nil {
		return "Password Hashing Error", fmt.Errorf("Failed to hash password: %v", err)
	}

	// ✅ FIX: Use the database client to store the user
	//NOTE: InsertUser requires a context parameter
	// ctx := context.Background()
	// err = dbClient.InsertUser(ctx, user)
	err = dbClient.InsertUser(context.Background(), user) // cleaner code
	if err != nil {
		return "Database Error", fmt.Errorf("Failed to create user: %v", err)
	}

	// Return success message
	return fmt.Sprintf("User %s created successfully!", event.Username), nil
}

//FIXME -
// Fixed Method Name: Changed CreateUser to InsertUser to match the actual method in the UserStore interface
// Added Context Parameter: Added context.Background() as the first parameter since InsertUser requires a context
// Added Context Import: Added the context package to the imports in main.go

// TODO - ProtectedHandler - to be moved to api/api.go
// ProtectedHandler is a handler function for protected routes that require authentication.
// This function demonstrates how to handle requests that have already passed through
// authentication middleware (like JWT validation).
//
// Parameters:
//   - request: APIGatewayProxyRequest containing HTTP request data from API Gateway
//
// Returns:
//   - APIGatewayProxyResponse: HTTP response that will be sent back to the client
//   - error: Any error that occurred during processing
//
// Usage: This handler is wrapped with ValidateJWTMiddleware to ensure only
// authenticated users can access this endpoint
func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Return a simple success response indicating the user has access to protected content
	// In a real application, this would contain actual protected data or perform
	// protected operations like accessing user-specific resources
	return events.APIGatewayProxyResponse{
		Body:       "This is a protected route", // Response body content
		StatusCode: http.StatusOK,               // HTTP 200 OK status code
	}, nil
}

// =============================================================
// !SECTION 3 - MAIN ENTRY POINT
// =============================================================
// main is the entry point of our AWS Lambda function.
// This function initializes the application and sets up the Lambda runtime
// to handle incoming API Gateway requests.
//
// How it works:
// 1. Creates a new application instance with all dependencies
// 2. Starts the Lambda runtime with a request handler function
// 3. The handler function routes requests based on the HTTP path
// 4. Each route is handled by appropriate handlers with middleware as needed
func main() {
	fmt.Println("====================================================================")
	fmt.Println("~~~~~~~~~~~~~~ Welcome to the AWS Lamb-da of Phoenix ! ~~~~~~~~~~~~~~")
	fmt.Println("=====================================================================")
	fmt.Println("................Running modes : Loco or Live ^_^ !!..................")
	fmt.Println("0 - Local testing with DOCKER for direct JSON (Backend Testing, n3rds)")
	fmt.Println("1 - AWS Lambda with direct JSON (HandleRequest-Frontend Testing, pimps)")
	fmt.Println("Default - Live on AWS Lambda with API Gateway (U, Reg. Joe - Sirs !)")

	// You can change this value to test different modes
	// Default: API Gateway - for production with middleware
	// Read from environment variable for flexibility
	//NOTE: Environment Variables are Always Strings so we can pass to Lambda as Global Environment Variable
	LocoOrLambdaStr := os.Getenv("LocoOrLambda")
	if LocoOrLambdaStr == "" {
		LocoOrLambdaStr = "99" // Default to API Gateway mode
	}

	//!TODO - LocoOrLambda Environment Variable Setting : THEKEY TO THE RUN MODE
	// Convert string to int
	LocoOrLambda := 99 // Default
	switch LocoOrLambdaStr {
	case "0":
		LocoOrLambda = 0
	case "1":
		LocoOrLambda = 1
	}

	switch LocoOrLambda {
	case 0: //LOCAL TESTING ONLY with DOCKER for JSON test input
		fmt.Println("Wassup, Loco!")

		// Local testing
		if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
			// Running locally
			testEvent := MyEvent{
				Username: "# Phonenix !",
				Password: "Password123 ?",
			}

			result, err := HandleRequest(testEvent)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Success: %s\n", result)
			}
			return
		}

		// AWS Lambda runtime
		lambda.Start(HandleRequest)

	case 1: //AWS Lambda with direct JSON
		fmt.Println("AWS Lambda with direct JSON (HandleRequest)")
		lambda.Start(HandleRequest)

	default: //DEFAULT: AWS Lambda with API Gateway
		fmt.Println("..... Running AWS Lambda with API Gateway (RegisterUser) .....")

		//FIXME: Lambda Function: Deploys and runs successfully, Database Operations: Connects to DynamoDB using environment variables
		//  Local Testing: Can run locally with go run main.go
		//  AWS Testing: Can test with AWS CLI or AWS Console
		// But To Make It Publicly Accessible: Need to use API Gateway >>> go to cdk5.go to uncomment the API Gateway integration

		// TODO: Uncomment this section when you want to use API Gateway integration
		// Initialize our application with all dependencies
		// This creates the dependency injection container that holds all our handlers
		// and services. In a real application, this would include database connections,
		// external service clients, configuration, etc.
		BTLambApp := app.NewApp()

		// Start the AWS Lambda runtime
		// This function tells AWS Lambda how to handle incoming requests
		// >>> This connects the Lambda function to the API Gateway <<<
		lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			// Route requests based on the HTTP path
			// This is a simple routing mechanism - in production, you might use
			// a more sophisticated router like gorilla/mux or gin
			switch request.Path {
			case "/register":
				// Handle user registration requests
				// This endpoint allows new users to create accounts
				// The handler will validate input, hash passwords, and store user data
				return BTLambApp.ApiHandler.RegisterUser(request)
				// takes the request from the API Gateway and passes it to the RegisterUser function in the ApiHandler layer

			case "/login":
				// Handle user login requests
				// This endpoint authenticates users and returns JWT tokens
				// The handler will validate credentials and generate authentication tokens
				return BTLambApp.ApiHandler.LoginUser(request)

			case "/protected":
				// Handle protected route requests
				// This endpoint requires authentication via JWT middleware
				// The middleware validates the JWT token before allowing access
				// This demonstrates the middleware pattern for authentication
				// wrap the ProtectedHandler with the ValidateJWTMiddleware in the middleware package
				return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)

			default:
				// Handle all other requests with a 404 Not Found response
				// This ensures that undefined routes return proper HTTP status codes
				return events.APIGatewayProxyResponse{
					Body:       "Got lost ? Not found", // Error message for undefined routes
					StatusCode: http.StatusNotFound,    // HTTP 404 Not Found status code
				}, nil
			}
		})
	}
}

//=============================================================
// !SECTION 4 - DOCUMENTATION & ANALYSIS
//=============================================================
/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA FUNCTION
===============================================================================

EXPLAINING THE CODE:
=========================
This main.go file serves as the entry point for a comprehensive AWS Lambda function
that demonstrates modern Go development practices, middleware patterns, and AWS integration.
It showcases how to build scalable, maintainable serverless applications.

• Entry Point: Main function that initializes the Lambda runtime
• Request Routing: Simple but effective path-based routing
• Middleware Integration: JWT authentication for protected routes
• Error Handling: Comprehensive error responses and logging
• Environment Configuration: Flexible deployment modes

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Standard Go packages (fmt, net/http, os)
   - AWS Lambda events package for API Gateway integration
   - Custom packages (app, database, middleware, types)

2. Data Structures:
   - MyEvent: Input data structure for direct Lambda testing
   - APIGatewayProxyRequest/Response: API Gateway integration types

3. Handler Functions:
   - HandleRequest: Direct Lambda function testing
   - ProtectedHandler: Protected route handler
   - main: Entry point and request routing

4. Integration Points:
   - API Gateway: HTTP request/response handling
   - Lambda Runtime: AWS Lambda function execution
   - Middleware: JWT authentication and validation

FUNCTION RELATIONSHIPS & CALL FLOW:
====================================
Request Processing Flow:
1. API Gateway → Lambda → main()
2. main() → Request Router → Handler Function
3. Handler → Middleware (if protected) → Business Logic
4. Business Logic → Database → Response
5. Response → API Gateway → Client

Authentication Flow:
1. Client → API Gateway → Lambda
2. Lambda → main() → Request Router
3. Router → middleware.ValidateJWTMiddleware
4. Middleware → JWT Validation → Handler
5. Handler → Protected Response

WHY WE USE THIS PATTERN:
========================
• Separation of Concerns: Each function has a single responsibility
• Middleware Pattern: Cross-cutting concerns like authentication
• Dependency Injection: Clean application structure
• Error Handling: Comprehensive error management
• AWS Integration: Native Lambda and API Gateway support

HOW IT WORKS:
=============
1. Application Initialization:
   - main() function starts the Lambda runtime
   - Environment variables determine execution mode
   - Application dependencies are injected

2. Request Processing:
   - API Gateway forwards HTTP requests to Lambda
   - Lambda function routes requests based on path
   - Middleware validates authentication for protected routes
   - Handlers process business logic and return responses

3. Response Handling:
   - Handlers return APIGatewayProxyResponse structures
   - Lambda runtime sends responses back to API Gateway
   - API Gateway forwards responses to clients

DEVELOPMENT MODES:
==================
This application supports three execution modes:

1. Local Testing (Mode 0):
   - Direct function invocation for testing
   - Uses local DynamoDB for development
   - No API Gateway integration

2. Direct Lambda (Mode 1):
   - AWS Lambda with direct JSON input
   - Uses AWS DynamoDB
   - No API Gateway integration

3. API Gateway (Mode 99 - Default):
   - Full AWS Lambda + API Gateway integration
   - HTTP request/response handling
   - Production-ready deployment

SECURITY FEATURES:
==================
• JWT Authentication: Secure token-based authentication
• Middleware Validation: Request validation before processing
• Error Handling: Secure error messages without information leakage
• Environment Configuration: Secure configuration management

PERFORMANCE CONSIDERATIONS:
===========================
• Cold Start Optimization: Efficient initialization
• Memory Management: Proper resource cleanup
• Database Connections: Connection pooling and reuse
• Response Caching: Optimized response handling

TESTING STRATEGY:
=================
• Unit Tests: Individual function testing
• Integration Tests: End-to-end request testing
• Local Testing: Docker-based development environment
• AWS Testing: Cloud-based testing and validation

DEPLOYMENT PROCESS:
===================
1. Code Development: Local development with Docker
2. Testing: Unit and integration testing
3. Build: Go compilation for Linux architecture
4. Package: ZIP file creation for Lambda deployment
5. Deploy: CDK deployment to AWS
6. Test: API Gateway and Lambda testing

CONCLUSION:
===========
This main.go file demonstrates a production-ready approach to AWS Lambda development
with Go. The modular structure, comprehensive error handling, and flexible deployment
modes make it suitable for both development and production environments.

Key Takeaways:
- Always implement proper error handling and logging
- Use middleware patterns for cross-cutting concerns
- Plan for multiple deployment modes (local, testing, production)
- Implement comprehensive request routing
- Follow Go best practices for maintainable code
- Use dependency injection for clean architecture
- Plan for security from the beginning
- Test thoroughly in all deployment modes

DEVELOPMENT ORDER REMINDER:
===========================
This is the LAST function to implement:
1. Types layer (types.go) - FIRST ✅
2. Database layer (database.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD ✅
4. API layer (api.go) - FOURTH ✅
5. App layer (app.go) - FIFTH ✅
6. Main entry (main.go) - LAST (THIS FILE) ✅

NEXT STEPS:
===========
1. Test all endpoints with different scenarios
2. Implement comprehensive logging and monitoring
3. Add performance optimization and caching
4. Set up CI/CD pipeline for automated deployment
5. Add comprehensive error handling and recovery
6. Implement security best practices
7. Add monitoring and alerting
8. Document API endpoints and usage
*/
