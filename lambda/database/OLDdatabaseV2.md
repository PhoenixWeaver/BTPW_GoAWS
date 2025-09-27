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
- DynamoDB integration with AWS SDK v2
- Repository pattern implementation
- Interface-based design for testability
- Error handling in database operations
- AWS session management

TOPICS COVERED:
- DynamoDB GetItem, PutItem operations
- AWS SDK v2 session management
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

Module Names Updated:
Root module: cdk4 → BT_GoAWS
Lambda module: lambda_func → BT_GoAWS_lambda
*/

package database

import (
	//  "lambda_func/types"
	"BTgoAWSlambda/types"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	// conflict with local types.go so we use the full path for types with alias dynamodbtypes
)

// ============================================================================
// CONSTANTS AND CONFIGURATION
// ============================================================================
// NOTE: Global variable TABLE_NAME is set by CDK to pass to Lambda as Global Environment Variable
// getTableName returns the DynamoDB table name from environment variables.
// This should match the table name defined in your CDK stack.
//
// Environment Variable: TABLE_NAME (set by CDK) - contains "BTtableGuestsInfo"
// Fallback: "BTtableGuestsInfo" for AWS production
//
// Security Note: This approach allows different table names for different environments
func getTableName() string {
	if tableName := os.Getenv("TABLE_NAME"); tableName != "" {
		// TABLE_NAME is set by CDK to pass to Lambda as Global Environment Variable
		// This returns "BTtableGuestsInfo" when running on AWS Lambda
		return tableName
	}
	// Fallback for AWS production
	return "BTtableGuestsInfo"
}

// getTableNameFromClient retrieves the table name from the client's tableName field
// This allows different clients to use different table names
// Returns: "BTtableGuestsInfo" for AWS or "LocoTable" for local
func (client *DynamoDBClient) getTableNameFromClient() string {
	if client.tableName != "" {
		return client.tableName // Returns "BTtableGuestsInfo" or "LocoTable"
	}
	// Fallback to environment variable ("BTtableGuestsInfo")
	return getTableName()
}

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
// NOTE: This interface is used to inject the database client into the API layer
// so that we can easily switch between different database implementations
// such as DynamoDB, PostgreSQL, etc.
// >>> in api.go, type ApiHandler struct {dbStore database.UserStore} instead of database.DynamoDBClient
type UserStore interface {
	// DoesUserExist checks if a user with the given username exists in the database.
	// Returns true if user exists, false if not, and error if database operation fails.
	DoesUserExist(ctx context.Context, username string) (bool, error)

	// InsertUser creates a new user in the database.
	// Returns error if user already exists or database operation fails.
	InsertUser(ctx context.Context, user types.User) error

	// GetUser retrieves a user from the database by username.
	// Returns user data if found, error if not found or database operation fails.
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
	databaseStore *dynamodb.Client // AWS DynamoDB service client (SDK v2)
	tableName     string           // DynamoDB table name: "BTtableGuestsInfo" or "LocoTable"
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

// /NOTE: previous version of dynamodb was using session.Must(session.NewSession()) + dynamodb.New(dbSession)
// / now it is using config.LoadDefaultConfig() + dynamodb.NewFromConfig(cfg)
func NewDynamoDB() DynamoDBClient {
	// Load default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO()) // context.TODO() is used here - in production >>>> app.go will pass in the proper context
	if err != nil {
		panic(err)
	}

	// Create DynamoDB client from configuration
	return DynamoDBClient{
		databaseStore: dynamodb.NewFromConfig(cfg),
	}
}

// STUB - Customized table name: local LocoTable_NAME or online BTtable_NAME
// NewDynamoDBWithTable creates a new DynamoDB client with a specific table name
func NewDynamoDBWithTable(tableName string) DynamoDBClient {
	// Load default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// Create DynamoDB client from configuration
	client := DynamoDBClient{
		databaseStore: dynamodb.NewFromConfig(cfg),
	}

	// Set the table name for this client
	client.tableName = tableName
	return client
}

// ============================================================================
// DATABASE OPERATIONS
// ============================================================================

// DoesUserExist checks if a user with the given username exists in DynamoDB.
// This function is used during user registration to prevent duplicate usernames.
//
// Parameters:
//   - ctx: Context for the operation
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
func (u DynamoDBClient) DoesUserExist(ctx context.Context, username string) (bool, error) {
	// Prepare GetItem request for DynamoDB
	result, err := u.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{ //GetItem retrieves a single item by its primary key
		TableName: aws.String(u.getTableNameFromClient()), //Table to query >>> NOTE: use client's table name
		Key: map[string]dynamodbtypes.AttributeValue{ //Primary key: username (partition key)
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: username, //String attribute value
			},
		},
	})

	// Handle database errors
	if err != nil {
		return true, fmt.Errorf("database error checking user existence: %w", err)
	} // Just to be safe, we return true and an error if the database error occurs (this will be handled by the calling code)

	// Check if item was found
	if result.Item == nil {
		// No item found - user doesn't really truly exist
		return false, nil
	}

	// Item found - user exists
	return true, nil
}

// InsertUser creates a new user in the DynamoDB table.
// This function is used during user registration to store new user data.
//
// Parameters:
//   - ctx: Context for the operation
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
func (u DynamoDBClient) InsertUser(ctx context.Context, user types.User) error {
	// Prepare PutItem request for DynamoDB
	item := &dynamodb.PutItemInput{ //PutItemInput is the input for the PutItem operation
		TableName: aws.String(u.getTableNameFromClient()), //Table to insert into >>> NOTE: use client's table name
		Item: map[string]dynamodbtypes.AttributeValue{
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: user.Username, // String attribute value >>> user.Username is the username from the user struct
			},
			"password": &dynamodbtypes.AttributeValueMemberS{
				Value: user.PasswordHash, // String attribute value >>> user.PasswordHash is the password hash from the user struct
			},
		},
	}

	// Execute the PutItem operation
	_, err := u.databaseStore.PutItem(ctx, item) //PutItem creates a new item or replaces an existing item
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Success - user inserted
	return nil
}

// GetUser retrieves a user from DynamoDB by username.
// This function is used during user login to retrieve user data for authentication.
//
// Parameters:
//   - ctx: Context for the operation
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
func (u DynamoDBClient) GetUser(ctx context.Context, username string) (types.User, error) {
	// Initialize empty user struct
	var user types.User

	// Prepare GetItem request for DynamoDB >>> same as DoesUserExist
	result, err := u.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{ //GetItem retrieves a single item by its primary key
		TableName: aws.String(u.getTableNameFromClient()), //Table to query >>> NOTE: use client's table name
		Key: map[string]dynamodbtypes.AttributeValue{ //Primary key: username (partition key)
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: username,
			},
		},
	})

	//STUB - need to check if the user exists
	// Handle database errors
	if err != nil {
		return user, fmt.Errorf("database error retrieving user: %w", err)
	}

	// Check if user was found
	if result.Item == nil {
		return user, fmt.Errorf("user not found: %s", username)
	}

	//STUB - need to deserialize the item to the user struct
	// Deserialize DynamoDB item to Go struct
	//  result.Item is the item from the DynamoDB table >>> need to unmarshal it to the user struct

	// DEBUG: Log the raw DynamoDB item before unmarshaling
	fmt.Printf("DEBUG: Raw DynamoDB item: %+v\n", result.Item)

	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, fmt.Errorf("failed to deserialize user data: %w", err)
	}

	// DEBUG: Log the user struct after unmarshaling
	fmt.Printf("DEBUG: User struct after unmarshaling: %+v\n", user)

	// Return successfully retrieved user
	return user, nil
}

//FIXME -
/*
What	                Name	               Value	                         Purpose
CDK  Constant	        TABLE_NAME	          "BTtableGuestsInfo"              	AWS table name
CDK Constant	        LOCAL_TABLE_NAME	  "LocoTable"	                    Local table name
Environment Variable	TABLE_NAME	          Set to TABLE_NAME value	        Passed to Lambda
*/
// ============================================================================
// COMPREHENSIVE CODE ANALYSIS
// ============================================================================

/*
EXPLAINING THE CODE:
=========================
This database package implements the repository pattern for user data management
using AWS DynamoDB SDK v2. It provides a clean interface for database operations and
follows Go best practices for error handling and type safety.

AWS SDK V2 MIGRATION CHANGES:
=============================
This implementation has been migrated from AWS SDK v1 to v2 with the following key changes:

1. Import Changes:
   - OLD: "github.com/aws/aws-sdk-go/aws/session"
   - NEW: "github.com/aws/aws-sdk-go-v2/config"
   - OLD: "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
   - NEW: "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

2. Client Creation:
   - OLD: session.Must(session.NewSession()) + dynamodb.New(dbSession)
   - NEW: config.LoadDefaultConfig() + dynamodb.NewFromConfig(cfg)

3. Type System Changes:
   - OLD: *dynamodb.AttributeValue (pointer-based)
   - NEW: dynamodbtypes.AttributeValue (interface-based)
   - OLD: map[string]*dynamodb.AttributeValue
   - NEW: map[string]dynamodbtypes.AttributeValue

4. Context Support:
   - OLD: No context parameter in methods
   - NEW: All methods require context.Context as first parameter

5. Attribute Value Creation:
   - OLD: &dynamodb.AttributeValue{S: aws.String(value)}
   - NEW: &dynamodbtypes.AttributeValueMemberS{Value: value}

6. Unmarshaling:
   - OLD: dynamodbattribute.UnmarshalMap()
   - NEW: attributevalue.UnmarshalMap()

7. Type Aliasing:
   - Added alias "dynamodbtypes" to avoid conflict with local "types" package

STRUCTURE OF THE CODE:
===========================
1. Interface Definition:
   - UserStore: Defines database operations contract with context support
   - Enables dependency injection and testing

2. Implementation:
   - DynamoDBClient: Implements UserStore interface using SDK v2
   - Handles AWS DynamoDB operations with proper context handling

3. Factory Function:
   - NewDynamoDB: Creates configured database client using config.LoadDefaultConfig()

4. CRUD Operations:
   - DoesUserExist: Check user existence with context
   - InsertUser: Create new user with context
   - GetUser: Retrieve user data with context

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
• AWS Integration: Native DynamoDB operations with SDK v2
• Type Safety: Go's type system prevents errors

AWS SDK V2 BENEFITS:
====================
• Context Support: Built-in support for cancellation and timeouts
• Better Performance: Improved connection pooling and reduced memory usage
• Type Safety: Stronger typing with interface-based attribute values
• Modular Design: Smaller, focused packages reduce bundle size
• Future-Proof: AWS's recommended SDK for new projects
• Better Error Handling: More descriptive error messages and types
• Simplified Configuration: Single config.LoadDefaultConfig() call

HOW IT WORKS:
=============
1. User Registration:
   - DoesUserExist(ctx, username) checks for duplicates with context
   - InsertUser(ctx, user) creates new user record with context

2. User Login:
   - GetUser(ctx, username) retrieves user data with context
   - Password validation happens in API layer

3. Context Handling:
   - All database operations require context.Context parameter
   - Context enables cancellation, timeouts, and request tracing
   - Context is passed through the entire call chain

4. Error Handling:
   - Database errors are wrapped with context information
   - Appropriate errors for different failure cases
   - Error wrapping preserves original error details

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
in AWS Lambda using the latest AWS SDK v2. The repository pattern implementation
makes it testable and maintainable, while proper error handling ensures robust operation.

Key Takeaways:
- Always use interfaces for database operations
- Implement proper error handling and wrapping
- Use AWS SDK v2 best practices with context support
- Design for testability from the start
- Handle edge cases (user not found, database errors)
- Use appropriate DynamoDB operations for each use case
- Implement proper logging and monitoring
- Consider performance implications of database operations
- Use type aliasing to avoid package name conflicts
- Leverage context for cancellation and timeout handling
- Take advantage of SDK v2's improved type safety
*/

// MockUserStore implements UserStore interface for testing
type MockUserStore struct {
	users map[string]types.User
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: make(map[string]types.User),
	}
}

func (m *MockUserStore) DoesUserExist(ctx context.Context, username string) (bool, error) {
	_, exists := m.users[username]
	return exists, nil
}

func (m *MockUserStore) InsertUser(ctx context.Context, user types.User) error {
	m.users[user.Username] = user
	return nil
}

func (m *MockUserStore) GetUser(ctx context.Context, username string) (types.User, error) {
	user, exists := m.users[username]
	if !exists {
		return types.User{}, fmt.Errorf("user not found: %s", username)
	}
	return user, nil
}