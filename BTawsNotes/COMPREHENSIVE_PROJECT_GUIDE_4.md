RequestID: e73f961d-2b75-4fcc-9e30-f19d4c358cbb

# AWS CDK + Lambda + DynamoDB Project - Comprehensive Recreation Guide

## 🎯 Project Overview

This is a complete AWS serverless application built with:
- **AWS CDK (Cloud Development Kit)** - Infrastructure as Code
- **AWS Lambda** - Serverless compute with Go runtime
- **AWS DynamoDB** - NoSQL database for user data
- **API Gateway** - REST API endpoints (currently commented out)
- **JWT Authentication** - Token-based authentication system

## 📋 Development Order - Step by Step

### **PHASE 1: Foundation Layer (Start Here)**

#### 1. **Database Layer** (`lambda/database/database.go`) - **FIRST FUNCTION**
```go
// Start with this function - it's the foundation
func NewDynamoDB() DynamoDBClient
```
**Why First:** Database operations are the foundation. All other layers depend on this.

**Key Functions to Implement:**
1. `NewDynamoDB()` - Database client creation
2. `DoesUserExist(ctx, username)` - Check user existence
3. `InsertUser(ctx, user)` - Create new user
4. `GetUser(ctx, username)` - Retrieve user data

**DynamoDB v2 Changes:**
- Uses `config.LoadDefaultConfig()` instead of session
- All methods require `context.Context` as first parameter
- Uses `dynamodbtypes.AttributeValue` instead of pointers
- Uses `attributevalue.UnmarshalMap()` for deserialization

#### 2. **Types Layer** (`lambda/types/types.go`) - **SECOND FUNCTION**
```go
// Second function - defines data structures
func NewUser(registerUser RegisterUser) (User, error)
```
**Why Second:** Defines data structures that database layer uses.

**Key Functions to Implement:**
1. `NewUser()` - Create user with hashed password
2. `ValidatePassword()` - Verify login credentials
3. `CreateToken()` - Generate JWT tokens

### **PHASE 2: Business Logic Layer**

#### 3. **Middleware Layer** (`lambda/middleware/middleware.go`) - **THIRD FUNCTION**
```go
// Third function - handles cross-cutting concerns
func ValidateJWTMiddleware(handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
```
**Why Third:** Provides authentication middleware for protected routes.

#### 4. **API Layer** (`lambda/api/api.go`) - **FOURTH FUNCTION**
```go
// Fourth function - handles HTTP requests
func (api *ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
func (api *ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
```
**Why Fourth:** Implements business logic and calls database layer.

**Key Functions to Implement:**
1. `NewApiHandler()` - Constructor with dependency injection
2. `RegisterUser()` - Handle user registration
3. `LoginUser()` - Handle user authentication

### **PHASE 3: Application Layer**

#### 5. **App Layer** (`lambda/app/app.go`) - **FIFTH FUNCTION**
```go
// Fifth function - dependency injection container
func NewApp() App
```
**Why Fifth:** Wires all dependencies together.

#### 6. **Main Entry** (`lambda/main.go`) - **LAST FUNCTION**
```go
// Last function - entry point and routing
func main()
```
**Why Last:** Orchestrates everything and handles request routing.

### **PHASE 4: Infrastructure Layer**

#### 7. **CDK Infrastructure** (`cdk4.go`) - **AFTER LAMBDA COMPLETE**
```go
// Infrastructure definition
func NewTake2GoCdkStack(scope constructs.Construct, id string, props *Take2GoCdkStackProps) awscdk.Stack
```
**Why Last:** Defines AWS resources after application logic is complete.

## 🔄 Function Call Flow & Relationships

### **Registration Flow:**
```
API Gateway → main() → RegisterUser() → DoesUserExist() → NewUser() → InsertUser()
```

### **Login Flow:**
```
API Gateway → main() → LoginUser() → GetUser() → ValidatePassword() → CreateToken()
```

### **Protected Route Flow:**
```
API Gateway → main() → ValidateJWTMiddleware() → ProtectedHandler()
```

## 🏗️ Project Structure

```
CDK4/
├── cdk4.go                    # CDK Infrastructure (Phase 4)
├── go.mod                     # Root module dependencies
├── lambda/
│   ├── main.go               # Entry point (Phase 3.2)
│   ├── app/
│   │   └── app.go           # Dependency injection (Phase 3.1)
│   ├── api/
│   │   └── api.go           # Business logic (Phase 2.2)
│   ├── middleware/
│   │   └── middleware.go    # Cross-cutting concerns (Phase 2.1)
│   ├── types/
│   │   └── types.go         # Data structures (Phase 1.2)
│   ├── database/
│   │   └── database.go      # Database operations (Phase 1.1)
│   └── go.mod               # Lambda module dependencies
```

## 🚀 Development Workflow

### **Step 1: Setup Environment**
```bash
# 1. Create project directory
mkdir CDK4
cd CDK4

# 2. Initialize Go modules
go mod init BT_GoAWS
cd lambda
go mod init BT_GoAWS_lambda
cd ..
```

### **Step 2: Implement Functions in Order**

#### **Phase 1.1: Database Layer**
```bash
# Create database/database.go
# Implement: NewDynamoDB, DoesUserExist, InsertUser, GetUser
```

#### **Phase 1.2: Types Layer**
```bash
# Create types/types.go
# Implement: NewUser, ValidatePassword, CreateToken
```

#### **Phase 2.1: Middleware Layer**
```bash
# Create middleware/middleware.go
# Implement: ValidateJWTMiddleware
```

#### **Phase 2.2: API Layer**
```bash
# Create api/api.go
# Implement: NewApiHandler, RegisterUser, LoginUser
```

#### **Phase 3.1: App Layer**
```bash
# Create app/app.go
# Implement: NewApp
```

#### **Phase 3.2: Main Entry**
```bash
# Update main.go
# Implement: main function with routing
```

#### **Phase 4: Infrastructure**
```bash
# Create cdk4.go
# Implement: CDK stack definition
```

### **Step 3: Build and Deploy**
```bash
# Build Lambda function
cd lambda
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bootstrap
Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

# Deploy infrastructure
cd ..
go run .
cdk deploy
```

## ⚠️ Common Mistakes to Avoid

### **1. Wrong Development Order**
- ❌ **Don't start with main.go** - You'll have nothing to wire together
- ❌ **Don't start with CDK** - You need the Lambda function first
- ✅ **Start with database layer** - It's the foundation everything else builds on

### **2. DynamoDB v2 Migration Issues**
- ❌ **Missing context parameter** - All DynamoDB v2 methods require context
- ❌ **Using old SDK imports** - Use `aws-sdk-go-v2` not `aws-sdk-go`
- ❌ **Wrong attribute value types** - Use `dynamodbtypes.AttributeValue` not pointers
- ✅ **Always pass context.Context** as first parameter
- ✅ **Use proper v2 imports and types**

### **3. Module Import Issues**
- ❌ **Wrong module names** - Must match go.mod module names
- ❌ **Missing imports** - Check all required packages are imported
- ✅ **Use consistent module naming** - BT_GoAWS for root, BT_GoAWS_lambda for lambda

### **4. Lambda Build Issues**
- ❌ **Wrong GOOS/GOARCH** - Must be linux/amd64 for AWS Lambda
- ❌ **Wrong handler name** - Must be "bootstrap" for Go runtime
- ❌ **Missing function.zip** - CDK expects this file
- ✅ **Always set environment variables** before building
- ✅ **Use correct file names** - bootstrap executable, function.zip package

### **5. CDK Deployment Issues**
- ❌ **Missing jsii.Close()** - Always defer this in main()
- ❌ **Wrong stack name** - Must be consistent across commands
- ❌ **Missing environment config** - Set account and region
- ✅ **Always defer jsii.Close()** in main function
- ✅ **Use consistent naming** throughout

## 🔧 Key Configuration Points

### **Module Names:**
- Root: `BT_GoAWS`
- Lambda: `BT_GoAWS_lambda`

### **Table Names:**
- DynamoDB: `userTable2-v2`

### **Function Names:**
- Lambda: `myLambdaFunction2`
- Stack: `Take2GoCdkStack`

### **AWS Configuration:**
- Account: `509505638450`
- Region: `us-east-1`

## 📚 Learning Progression

### **Beginner Level:**
1. Understand Go basics and package structure
2. Learn AWS Lambda concepts
3. Understand DynamoDB operations

### **Intermediate Level:**
1. Implement repository pattern
2. Understand dependency injection
3. Learn JWT authentication

### **Advanced Level:**
1. Master CDK infrastructure patterns
2. Implement proper error handling
3. Add monitoring and logging

## 🎯 Next Steps for Improvement

### **Immediate Improvements:**
1. Uncomment API Gateway in CDK
2. Add IAM permissions
3. Implement proper error handling
4. Add input validation

### **Production Ready:**
1. Add environment variables
2. Implement proper logging
3. Add monitoring and alarms
4. Implement rate limiting
5. Add security headers

### **Advanced Features:**
1. Add user profile management
2. Implement password reset
3. Add email verification
4. Implement role-based access control
5. Add API versioning

## 🔍 Testing Strategy

### **Unit Testing:**
- Test each function independently
- Mock database operations
- Test error conditions

### **Integration Testing:**
- Test complete user flows
- Test API endpoints
- Test database operations

### **End-to-End Testing:**
- Test complete deployment
- Test real AWS resources
- Test performance under load

## 📖 Additional Resources

### **AWS Documentation:**
- [AWS CDK Go](https://docs.aws.amazon.com/cdk/v2/guide/work-with-cdk-go.html)
- [AWS Lambda Go](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)
- [DynamoDB Go SDK v2](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/dynamodb-example.html)

### **Go Best Practices:**
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

This guide provides a complete roadmap for recreating this project from scratch, following the proper development order and avoiding common pitfalls.
