/*
===============================================================================
AWS LAMBDA FUNCTION - APPLICATION LAYER (DEPENDENCY INJECTION CONTAINER)
===============================================================================

DEVELOPMENT ORDER: FIFTH FUNCTION TO IMPLEMENT
==============================================
This is the FIFTH function in the development sequence:
1. Types layer (types.go) - FIRST ✅
2. Database layer (database.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD ✅
4. API layer (api.go) - FOURTH ✅
5. App layer (app.go) - FIFTH (THIS FILE) ✅
6. Main entry (main.go) - LAST

Author: Ben Tran
Date: 22/09/2025
Description: This is the APPLICATION LAYER that orchestrates all dependencies using
             the Dependency Injection pattern. It creates and wires together all
             the components needed for the Lambda function to work.

LEARNING OBJECTIVES:
- Dependency Injection pattern implementation
- Application layer architecture
- Component wiring and initialization
- Factory pattern for application creation
- Separation of concerns between layers

TOPICS COVERED:
- Dependency injection container pattern
- Factory function implementation
- Component lifecycle management
- Interface-based dependency management
- Application initialization and configuration

FEATURES:
- Centralized dependency management
- Clean separation of concerns
- Easy testing through dependency injection
- Modular application architecture
- Single responsibility principle

ARCHITECTURE:
- App struct: Main application container
- NewApp(): Factory function for application creation
- Dependency injection: Database → API Handler → App
- Interface-based design: Uses UserStore interface
- AWS SDK v2 integration: DynamoDB client creation

FUNCTION RELATIONSHIPS & DEPENDENCIES:
======================================
1. NewApp() → NewDynamoDB() → DynamoDBClient (Database layer)
2. NewApp() → NewApiHandler() → ApiHandler (API layer)
3. App struct → ApiHandler (Contains API handler)
4. ApiHandler → UserStore interface (Database operations)

Layer Dependencies (Bottom-up):
- App Layer (THIS) → API Layer (api.go)
- App Layer (THIS) → Database Layer (database.go)
- Main Layer (main.go) → App Layer (THIS)

DYNAMODB V2 MIGRATION NOTES:
============================
This file has been updated to work with AWS SDK v2:
- NewDynamoDB() now uses config.LoadDefaultConfig() instead of session.NewSession()
- DynamoDBClient implements UserStore interface with context support
- All database operations now require context.Context parameter
- Uses new v2 client patterns: dynamodb.NewFromConfig(cfg)
- Interface-based design allows for easy testing and mocking

INTEGRATION POINTS:
===================
- Uses database.UserStore interface for database operations
- Uses api.ApiHandler for request processing
- Returns App struct for main.go to use
- Handles dependency injection and component wiring
*/

// Package app implements the application layer using the Dependency Injection pattern.
// This package is responsible for creating and wiring together all the components
// needed for the Lambda function to work properly.
package app

import (
	"BTgoAWSlambda/api"      // API layer for request handling
	"BTgoAWSlambda/database" // Database layer for data persistence
)

// KEY POINT: Import paths use the module name "BTgoAWSlambda" defined in go.mod
// This ensures proper module resolution and avoids import path issues

// App represents the main application container that holds all the components
// needed for the Lambda function to operate. This struct follows the Dependency
// Injection pattern, making the application testable and maintainable.
//
// Fields:
//   - ApiHandler: The API handler that processes HTTP requests from API Gateway
//
// Usage: Created by NewApp(), used by main.go for request routing
// Example: app := NewApp(); response := app.ApiHandler.RegisterUser(request)
type App struct {
	ApiHandler api.ApiHandler // API handler with all business logic
}

// NewApp creates a new App instance with all dependencies properly wired together.
// This function implements the Factory pattern and Dependency Injection pattern,
// creating all necessary components and injecting their dependencies.
//
// DEPENDENCY INJECTION FLOW:
// ==========================
// 1. Create database client (DynamoDB v2)
// 2. Inject database client into API handler
// 3. Return App struct with configured API handler
//
// DYNAMODB V2 CHANGES:
// ====================
// - Uses config.LoadDefaultConfig() instead of session.NewSession()
// - Creates DynamoDBClient that implements UserStore interface
// - All database operations now support context.Context
// - Uses new v2 client patterns: dynamodb.NewFromConfig(cfg)
//
// Returns:
//   - App: Configured application ready for use
//
// Usage: Called during Lambda function initialization in main.go
// Example: lambdaApp := app.NewApp()
//
// CRITICAL SUCCESS FACTORS:
// =========================
// 1. Database client must implement UserStore interface ✅
// 2. API handler must accept UserStore interface ✅
// 3. All dependencies must be properly injected ✅
// 4. Module imports must use correct paths ✅
func NewApp() App {
	// STEP 1: Create database client using AWS SDK v2
	// KEY POINT: NewDynamoDB() returns DynamoDBClient which implements UserStore interface
	// This allows for easy testing by mocking the UserStore interface
	db := database.NewDynamoDB() // Creates DynamoDB client with v2 SDK
	// ✅ >>>>> TO SWITCH BETWEEN DATABASES: ProgreSQL Mongo etc >>>>>> RIGHT HERE

	// STEP 2: Create API handler with database dependency injection
	// KEY POINT: NewApiHandler() expects UserStore interface, not concrete type
	// This follows the Dependency Inversion Principle (DIP)
	apiHandler := api.NewApiHandler(db) // Injects database dependency

	// STEP 3: Return configured application
	// KEY POINT: App struct contains all necessary components for Lambda function
	return App{
		ApiHandler: apiHandler, // API handler ready to process requests
	}
}

/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS LAMBDA APPLICATION LAYER
===============================================================================

EXPLAINING THE CODE:
=========================
This application layer implements the Dependency Injection pattern for AWS Lambda functions
with the following architecture:

• Dependency Flow: Database → API Handler → App → Main
• Interface-Based Design: Uses UserStore interface for database operations
• Factory Pattern: NewApp() creates and configures all components
• AWS SDK v2 Integration: DynamoDB client with context support

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - Uses module name "BTgoAWSlambda" for proper import resolution
   - Imports api and database packages for dependency injection

2. Data Structures:
   - App struct: Main application container with ApiHandler field
   - Follows single responsibility principle

3. Factory Functions:
   - NewApp(): Creates and wires all dependencies together
   - Implements dependency injection pattern

4. Dependency Relationships:
   - App → ApiHandler → UserStore interface → DynamoDBClient

FUNCTION RELATIONSHIPS & CALL FLOW:
===================================
1. NewApp() → NewDynamoDB() → DynamoDBClient (Database creation)
2. NewApp() → NewApiHandler() → ApiHandler (API handler creation)
3. App struct → ApiHandler (Contains configured handler)
4. ApiHandler → UserStore interface (Database operations)

Layer Dependencies (Bottom-up):
- App Layer (THIS) → API Layer (api.go)
- App Layer (THIS) → Database Layer (database.go)
- Main Layer (main.go) → App Layer (THIS)

WHY WE USE THIS PATTERN:
========================
• Dependency Injection: Makes testing easier and code more flexible
• Interface-Based Design: Allows for easy mocking and testing
• Single Responsibility: Each layer has a specific purpose
• Factory Pattern: Centralizes object creation and configuration
• Clean Architecture: Separates concerns between layers

HOW IT WORKS:
=============
1. Lambda function starts and calls NewApp()
2. NewApp() creates DynamoDB client using AWS SDK v2
3. NewApp() injects database client into API handler
4. NewApp() returns configured App struct
5. Main.go uses App.ApiHandler for request processing
6. API handler uses injected database for operations

DYNAMODB V2 MIGRATION BENEFITS:
===============================
• Context Support: All operations support cancellation and timeouts
• Better Performance: Improved connection pooling and reduced memory usage
• Type Safety: Stronger typing with interface-based attribute values
• Future-Proof: AWS's recommended SDK for new projects
• Simplified Configuration: Single config.LoadDefaultConfig() call

APPLICATION IN USE:
==================
• Lambda Initialization: Called once when Lambda starts
• Request Processing: App.ApiHandler processes all requests
• Database Operations: All database calls go through UserStore interface
• Error Handling: Proper error propagation through all layers

IMPROVEMENT IDEAS:
==================
• Add configuration management for environment variables
• Implement proper context with timeouts and cancellation
• Add logging and monitoring for application lifecycle
• Implement health check endpoints
• Add metrics collection for performance monitoring
• Use AWS Secrets Manager for sensitive configuration
• Implement graceful shutdown handling
• Add request/response logging middleware

SECURITY CONSIDERATIONS:
========================
• Interface-Based Design: Prevents direct database access
• Dependency Injection: Allows for secure testing
• AWS SDK v2: Uses latest security features
• Context Support: Enables request tracing and cancellation
• Error Handling: Prevents information leakage

PERFORMANCE OPTIMIZATIONS:
==========================
• Single Instance: App created once per Lambda container
• Connection Pooling: Handled by AWS SDK v2
• Interface-Based: Minimal overhead for dependency injection
• Lazy Loading: Components created only when needed
• Memory Efficiency: Shared instances across requests

TESTING STRATEGY:
=================
• Mock UserStore: Create mock implementation for testing
• Dependency Injection: Easy to inject test dependencies
• Interface-Based: Test API handler without real database
• Unit Tests: Test each component independently
• Integration Tests: Test complete application flow

CONCLUSION:
===========
This application layer provides a solid foundation for AWS Lambda functions
using modern Go patterns and AWS SDK v2. The dependency injection pattern
makes the code testable, maintainable, and follows clean architecture principles.

Key Takeaways:
- Always use dependency injection for testability
- Follow the exact development order: Types → Database → Middleware → API → App → Main
- Use interface-based design for flexibility
- Implement proper error handling and logging
- Use AWS SDK v2 for better performance and features
- Keep components loosely coupled and highly cohesive
- Plan for testing from the beginning
- Use factory patterns for object creation
- Follow single responsibility principle
- Implement proper context handling for AWS SDK v2

DEVELOPMENT ORDER REMINDER:
===========================
This is the FIFTH function to implement:
1. Types layer (types.go) - FIRST ✅
2. Database layer (database.go) - SECOND ✅
3. Middleware layer (middleware.go) - THIRD ✅
4. API layer (api.go) - FOURTH ✅
5. App layer (app.go) - FIFTH (THIS FILE) ✅
6. Main entry (main.go) - LAST

NEXT STEPS:
===========
1. Implement main.go with request routing
2. Set up CDK infrastructure
3. Add comprehensive testing
4. Implement monitoring and logging
5. Add environment configuration
6. Deploy and test the complete application
*/
