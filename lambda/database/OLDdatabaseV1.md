/*
===============================================================================
AWS LAMBDA FUNCTION - DATABASE PACKAGE
===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This package contains all database operations for user management
             using AWS DynamoDB. It implements the repository pattern for
             clean separation of concerns and testability.

LEARNING OBJECTIVES:
- DynamoDB integration with AWS SDK
- Repository pattern implementation
- Interface-based design for testability
- Error handling in database operations
- AWS session management

TOPICS COVERED:
- DynamoDB GetItem, PutItem operations
- AWS SDK session management
- Interface design for dependency injection
- Error handling and validation
- Data serialization/deserialization

FEATURES:
- User existence checking
- User insertion with validation
- User retrieval by username
- Interface-based design for testing
- Proper error handling

ARCHITECTURE:
- UserStore interface: Defines database operations
- DynamoDBClient: Implements UserStore interface
- NewDynamoDB: Factory function for database client
- CRUD operations: Create, Read, Update, Delete

FUNCTION RELATIONSHIPS:
1. NewDynamoDB → DynamoDBClient (Database client creation)
2. DoesUserExist → GetItem (Check if user exists)
3. InsertUser → PutItem (Create new user)
4. GetUser → GetItem (Retrieve user data)

DEVELOPMENT ORDER:
1. Database layer (THIS FILE) - START HERE
2. Types layer (types.go) - SECOND
3. Middleware layer (middleware.go) - THIRD
4. API layer (api.go) - FOURTH
5. App layer (app.go) - FIFTH
6. Main entry (main.go) - LAST

*/

package database

// import (
// 	"fmt"
// 	"lambda-func/types"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// )

///FIXME - new version of dynamodb
import (
	"context"
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ============================================================================
// CONSTANTS AND CONFIGURATION
// ============================================================================

// TABLE_NAME is the DynamoDB table name for user data.
// This should match the table name defined in your CDK stack.
//
// Security Note: In production, consider using environment variables
// for table names to support different environments (dev, staging, prod)
const TABLE_NAME = "userTable2"

// ============================================================================
// INTERFACE DEFINITIONS
// ============================================================================

// UserStore defines the interface for all user-related database operations.
// This interface allows for easy testing and dependency injection.
//
// Interface Benefits:
// - Testability: Can create mock implementations for testing
// - Flexibility: Can switch between different database implementations
// - Clean Architecture: Separates business logic from data access
//
// Usage: Implemented by DynamoDBClient, can be mocked for testing
type UserStore interface {
	// // DoesUserExist checks if a user with the given username exists in the database.
	// // Returns true if user exists, false if not, and error if database operation fails.
	// DoesUserExist(username string) (bool, error)

	// // InsertUser creates a new user in the database.
	// // Returns error if user already exists or database operation fails.
	// InsertUser(user types.User) error

	// // GetUser retrieves a user from the database by username.
	// // Returns user data if found, error if not found or database operation fails.
	// GetUser(username string) (types.User, error)

	///FIXME - for new version of dynamodb, we need to use the context.Context

	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user types.User) error
	GetUser(ctx context.Context, username string) (types.User, error)
}

// ============================================================================
// DATABASE CLIENT IMPLEMENTATION
// ============================================================================

// DynamoDBClient implements the UserStore interface using AWS DynamoDB.
// This struct contains the DynamoDB service client for database operations.
//
// Fields:
// - databaseStore: AWS DynamoDB service client
//
// Usage: Created by NewDynamoDB(), used for all database operations
type DynamoDBClient struct {
	// databaseStore *dynamodb.DynamoDB // AWS DynamoDB service client
	databaseStore *dynamodb.Client /// new version of dynamodb
}

// ============================================================================
// FACTORY FUNCTIONS
// ============================================================================

// NewDynamoDB creates a new DynamoDBClient with proper AWS session configuration.
// This function initializes the AWS session and DynamoDB service client.
//
// Returns:
//   - DynamoDBClient: Configured database client ready for use
//
// AWS Configuration:
// - Uses default AWS credentials (from environment, IAM role, or AWS config)
// - Uses default region (from environment or AWS config)
// - Creates new session for each client instance
//
// Usage: Called during application initialization to create database client
func NewDynamoDB() DynamoDBClient {
	// // Create AWS session with default configuration
	// // session.Must() panics if session creation fails
	// // In production, you might want to handle this error gracefully
	// dbSession := session.Must(session.NewSession())

	// // Create DynamoDB service client from session
	// db := dynamodb.New(dbSession)

	// // Return configured client
	// return DynamoDBClient{
	// 	databaseStore: db,
	// }

	///FIXME - new version of dynamodb need to use the context.Context

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return DynamoDBClient{
		databaseStore: dynamodb.NewFromConfig(cfg),
	}
}

// ============================================================================
// DATABASE OPERATIONS
// ============================================================================

// DoesUserExist checks if a user with the given username exists in DynamoDB.
// This function is used during user registration to prevent duplicate usernames.
//
// Parameters:
//   - username: The username to check for existence
//
// Returns:
//   - bool: true if user exists, false if not
//   - error: Any error that occurred during database operation
//
// DynamoDB Operation:
// - Uses GetItem to check for user existence
// - Returns true if item exists, false if not
// - Handles database errors appropriately
//
// Usage: Called during user registration to validate username uniqueness
// func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
// Prepare GetItem request for DynamoDB
// GetItem retrieves a single item by its primary key
//
//	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
//		TableName: aws.String(TABLE_NAME), // Table to query
//		Key: map[string]*dynamodb.AttributeValue{
//			// Primary key: username (partition key)
//			"username": {
//				S: aws.String(username), // String attribute value
//			},
//		},
//	})
//
// /FIXME - new version of dynamodb
func (u DynamoDBClient) DoesUserExist(ctx context.Context, username string) (bool, error) {
	result, err := u.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{
				Value: username,
			},
		},
	})
	// Handle database errors
	if err != nil {
		// Return true and error - this indicates a database problem
		// The calling code should handle this error appropriately
		return true, fmt.Errorf("database error checking user existence: %w", err)
	}

	// Check if item was found
	if result.Item == nil {
		// No item found - user doesn't exist
		return false, nil
	}

	// Item found - user exists
	return true, nil
}

// InsertUser creates a new user in the DynamoDB table.
// This function is used during user registration to store new user data.
//
// Parameters:
//   - user: The user data to insert (must have Username and PasswordHash)
//
// Returns:
//   - error: Any error that occurred during database operation
//
// DynamoDB Operation:
// - Uses PutItem to create new user record
// - Stores username as partition key and password hash
// - Overwrites existing item if username already exists
//
// Security Note: PasswordHash should already be hashed before calling this function
//
// Usage: Called during user registration after validation
// func (u DynamoDBClient) InsertUser(user types.User) error {
// Prepare PutItem request for DynamoDB
// PutItem creates a new item or replaces an existing item
//
//	item := &dynamodb.PutItemInput{
//		TableName: aws.String(TABLE_NAME), // Table to insert into
//		Item: map[string]*dynamodb.AttributeValue{
//			// Username as partition key (primary key)
//			"username": {
//				S: aws.String(user.Username), // String attribute value
//			},
//			// Password hash (should be pre-hashed)
//			"password": {
//				S: aws.String(user.PasswordHash), // String attribute value
//			},
//		},
//	}
//
// /FIXME - new version of dynamodb
func (u DynamoDBClient) InsertUser(ctx context.Context, user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{
				Value: user.Username,
			},
			"password": &types.AttributeValueMemberS{
				Value: user.PasswordHash,
			},
		},
	}

	// Execute the PutItem operation
	_, err := u.databaseStore.PutItem(ctx, item)
	if err != nil {
		// Return error if database operation fails
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Success - user inserted
	return nil
}

// GetUser retrieves a user from DynamoDB by username.
// This function is used during user login to retrieve user data for authentication.
//
// Parameters:
//   - username: The username to retrieve
//
// Returns:
//   - types.User: The user data if found
//   - error: Error if user not found or database operation fails
//
// DynamoDB Operation:
// - Uses GetItem to retrieve user by primary key
// - Deserializes DynamoDB item to Go struct
// - Returns appropriate errors for different failure cases
//
// Usage: Called during user login to retrieve user data for password validation
// func (u DynamoDBClient) GetUser(username string) (types.User, error) {
// 	// Initialize empty user struct
// 	var user types.User

// 	// Prepare GetItem request for DynamoDB
// 	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
// 		TableName: aws.String(TABLE_NAME), // Table to query
// 		Key: map[string]*dynamodb.AttributeValue{
// 			// Primary key: username (partition key)
// 			"username": {
// 				S: aws.String(username), // String attribute value
// 			},
// 		},
// 	})

// /FIXME - new version of dynamodb
func (u DynamoDBClient) GetUser(ctx context.Context, username string) (types.User, error) {
	// Initialize empty user struct
	var user types.User
	result, err := u.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{
				Value: username,
			},
		},
	})
	// Handle database errors
	if err != nil {
		return user, fmt.Errorf("database error retrieving user: %w", err)
	}

	// Check if user was found
	if result.Item == nil {
		// User not found - return appropriate error
		return user, fmt.Errorf("user not found: %s", username)
	}

	// Deserialize DynamoDB item to Go struct
	// dynamodbattribute.UnmarshalMap converts DynamoDB attributes to Go struct
	// err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	///FIXME - new version of dynamodb
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, fmt.Errorf("failed to deserialize user data: %w", err)
	}

	// Return successfully retrieved user
	return user, nil
}

// ============================================================================
// COMPREHENSIVE CODE ANALYSIS
// ============================================================================

/*
EXPLAINING THE CODE:
=========================
This database package implements the repository pattern for user data management
using AWS DynamoDB. It provides a clean interface for database operations and
follows Go best practices for error handling and type safety.

STRUCTURE OF THE CODE:
===========================
1. Interface Definition:
   - UserStore: Defines database operations contract
   - Enables dependency injection and testing

2. Implementation:
   - DynamoDBClient: Implements UserStore interface
   - Handles AWS DynamoDB operations

3. Factory Function:
   - NewDynamoDB: Creates configured database client

4. CRUD Operations:
   - DoesUserExist: Check user existence
   - InsertUser: Create new user
   - GetUser: Retrieve user data

FUNCTION RELATIONSHIPS:
========================
Database Operations Flow:
1. NewDynamoDB → DynamoDBClient (Client creation)
2. DoesUserExist → GetItem (Existence check)
3. InsertUser → PutItem (User creation)
4. GetUser → GetItem (User retrieval)

Integration with Other Layers:
1. Database → Types (Uses types.User struct)
2. Database → API (Called by API handlers)
3. Database → App (Injected as dependency)

WHY WE USE THIS PATTERN:
========================
• Repository Pattern: Clean separation of data access logic
• Interface Design: Enables testing and flexibility
• Error Handling: Proper error propagation and handling
• AWS Integration: Native DynamoDB operations
• Type Safety: Go's type system prevents errors

HOW IT WORKS:
=============
1. User Registration:
   - DoesUserExist checks for duplicates
   - InsertUser creates new user record

2. User Login:
   - GetUser retrieves user data
   - Password validation happens in API layer

3. Error Handling:
   - Database errors are wrapped with context
   - Appropriate errors for different failure cases

SECURITY CONSIDERATIONS:
========================
• Password Hashing: Passwords should be hashed before storage
• Input Validation: Username validation should happen in API layer
• Error Messages: Don't expose internal database errors
• Access Control: Use IAM roles for database access
• Encryption: DynamoDB encryption at rest and in transit

IMPROVEMENT IDEAS:
==================
• Add connection pooling for better performance
• Implement retry logic for transient failures
• Add metrics and logging for monitoring
• Use environment variables for table names
• Implement batch operations for multiple users
• Add data validation before database operations
• Use AWS X-Ray for tracing
• Implement caching layer for frequently accessed data

PERFORMANCE OPTIMIZATIONS:
==========================
• Use DynamoDB batch operations for multiple items
• Implement connection pooling
• Add caching layer for frequently accessed data
• Use DynamoDB Global Secondary Indexes (GSI) for queries
• Optimize item size and attribute access
• Use DynamoDB streams for real-time updates

CONCLUSION:
===========
This database package provides a solid foundation for user data management
in AWS Lambda. The repository pattern implementation makes it testable and
maintainable, while proper error handling ensures robust operation.

Key Takeaways:
- Always use interfaces for database operations
- Implement proper error handling and wrapping
- Use AWS SDK best practices
- Design for testability from the start
- Handle edge cases (user not found, database errors)
- Use appropriate DynamoDB operations for each use case
- Implement proper logging and monitoring
- Consider performance implications of database operations
*/
