# AWS Lambda Go Development Guide
## Complete Step-by-Step Development Process

### Table of Contents
1. [Project Setup & Structure](#project-setup--structure)
2. [Development Order & Function Priority](#development-order--function-priority)
3. [Step-by-Step Implementation](#step-by-step-implementation)
4. [Best Practices & Common Pitfalls](#best-practices--common-pitfalls)
5. [Testing Strategy](#testing-strategy)
6. [Deployment Considerations](#deployment-considerations)

---

## Project Setup & Structure

### 1. Initial Project Structure
```
lambda/
├── go.mod                 # Module definition
├── go.sum                 # Dependency checksums
├── main.go               # Entry point
├── app/
│   └── app.go           # Application logic
├── middleware/
│   └── middleware.go    # Middleware functions
├── handlers/
│   └── handlers.go      # Request handlers
├── models/
│   └── models.go        # Data structures
└── utils/
    └── utils.go         # Utility functions
```

### 2. Module Initialization
```bash
# Initialize Go module
go mod init lambda_func

# Add AWS Lambda dependencies
go get github.com/aws/aws-lambda-go/lambda
go get github.com/aws/aws-lambda-go/events
```

---

## Development Order & Function Priority

### **CRITICAL: Development Sequence**
Follow this exact order to avoid dependency issues:

1. **Models/Data Structures** (First)
2. **Utility Functions** (Second)
3. **Middleware** (Third)
4. **Handlers** (Fourth)
5. **App Structure** (Fifth)
6. **Main Function** (Last)

---

## Step-by-Step Implementation

### Step 1: Create Models (models/models.go)
```go
package models

import "time"

// User represents a user in the system
type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginRequest represents user login data
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// APIResponse represents standard API response
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

### Step 2: Create Utility Functions (utils/utils.go)
```go
package utils

import (
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
    "net/http"
)

// CreateSuccessResponse creates a successful API response
func CreateSuccessResponse(data interface{}, message string) events.APIGatewayProxyResponse {
    response := APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    }
    
    body, _ := json.Marshal(response)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(body),
    }
}

// CreateErrorResponse creates an error API response
func CreateErrorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
    response := APIResponse{
        Success: false,
        Error:   message,
    }
    
    body, _ := json.Marshal(response)
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(body),
    }
}

// ParseJSONBody parses JSON from request body
func ParseJSONBody(body string, target interface{}) error {
    return json.Unmarshal([]byte(body), target)
}
```

### Step 3: Create Middleware (middleware/middleware.go)
```go
package middleware

import (
    "github.com/aws/aws-lambda-go/events"
    "lambda_func/utils"
)

// ValidateJWTMiddleware validates JWT tokens
func ValidateJWTMiddleware(handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        // Extract token from Authorization header
        authHeader := request.Headers["Authorization"]
        if authHeader == "" {
            return utils.CreateErrorResponse("Authorization header required", 401), nil
        }
        
        // TODO: Implement JWT validation logic
        // For now, just pass through
        return handler(request)
    }
}

// LoggingMiddleware logs requests
func LoggingMiddleware(handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        // Log request
        fmt.Printf("Request: %s %s\n", request.HTTPMethod, request.Path)
        
        // Call handler
        response, err := handler(request)
        
        // Log response
        fmt.Printf("Response: %d\n", response.StatusCode)
        
        return response, err
    }
}
```

### Step 4: Create Handlers (handlers/handlers.go)
```go
package handlers

import (
    "github.com/aws/aws-lambda-go/events"
    "lambda_func/models"
    "lambda_func/utils"
)

type UserHandler struct {
    // Add dependencies here (database, services, etc.)
}

func NewUserHandler() *UserHandler {
    return &UserHandler{}
}

func (h *UserHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var req models.RegisterRequest
    
    // Parse request body
    if err := utils.ParseJSONBody(request.Body, &req); err != nil {
        return utils.CreateErrorResponse("Invalid JSON", 400), nil
    }
    
    // Validate required fields
    if req.Username == "" || req.Email == "" || req.Password == "" {
        return utils.CreateErrorResponse("Username, email, and password are required", 400), nil
    }
    
    // TODO: Implement user registration logic
    // - Hash password
    // - Save to database
    // - Generate JWT token
    
    user := models.User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    return utils.CreateSuccessResponse(user, "User registered successfully"), nil
}

func (h *UserHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var req models.LoginRequest
    
    // Parse request body
    if err := utils.ParseJSONBody(request.Body, &req); err != nil {
        return utils.CreateErrorResponse("Invalid JSON", 400), nil
    }
    
    // Validate required fields
    if req.Username == "" || req.Password == "" {
        return utils.CreateErrorResponse("Username and password are required", 400), nil
    }
    
    // TODO: Implement login logic
    // - Validate credentials
    // - Generate JWT token
    
    return utils.CreateSuccessResponse(map[string]string{
        "token": "jwt_token_here",
    }, "Login successful"), nil
}
```

### Step 5: Create App Structure (app/app.go)
```go
package app

import "lambda_func/handlers"

type App struct {
    UserHandler *handlers.UserHandler
}

func NewApp() *App {
    return &App{
        UserHandler: handlers.NewUserHandler(),
    }
}
```

### Step 6: Create Main Function (main.go)
```go
package main

import (
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "lambda_func/app"
    "lambda_func/middleware"
    "lambda_func/utils"
    "net/http"
)

func main() {
    lambdaApp := app.NewApp()
    
    lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        // Apply logging middleware to all requests
        handler := middleware.LoggingMiddleware(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
            switch req.Path {
            case "/register":
                return lambdaApp.UserHandler.RegisterUser(req)
            case "/login":
                return lambdaApp.UserHandler.LoginUser(req)
            case "/protected":
                // Apply JWT middleware to protected routes
                protectedHandler := middleware.ValidateJWTMiddleware(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
                    return utils.CreateSuccessResponse(nil, "This is a protected route"), nil
                })
                return protectedHandler(req)
            default:
                return utils.CreateErrorResponse("Not found", http.StatusNotFound), nil
            }
        })
        
        return handler(request)
    })
}
```

---

## Best Practices & Common Pitfalls

### ✅ DO's

1. **Always start with models/data structures**
   - Define your data contracts first
   - Use proper JSON tags
   - Include validation tags

2. **Use proper error handling**
   - Return meaningful error messages
   - Use appropriate HTTP status codes
   - Log errors for debugging

3. **Implement middleware pattern**
   - Separate concerns (logging, auth, validation)
   - Make handlers composable
   - Keep handlers focused on business logic

4. **Use dependency injection**
   - Pass dependencies to handlers
   - Make testing easier
   - Improve code maintainability

5. **Follow REST conventions**
   - Use proper HTTP methods
   - Return consistent response formats
   - Use appropriate status codes

### ❌ DON'Ts

1. **Don't start with main.go**
   - You'll hit import errors
   - Hard to test individual components
   - Poor separation of concerns

2. **Don't ignore error handling**
   - Always check for errors
   - Don't panic in production
   - Log errors appropriately

3. **Don't hardcode values**
   - Use environment variables
   - Use configuration files
   - Make your code configurable

4. **Don't mix concerns**
   - Keep handlers focused
   - Separate business logic from HTTP handling
   - Use middleware for cross-cutting concerns

---

## Testing Strategy

### Unit Testing
```go
// handlers/handlers_test.go
package handlers

import (
    "testing"
    "github.com/aws/aws-lambda-go/events"
)

func TestRegisterUser(t *testing.T) {
    handler := NewUserHandler()
    
    request := events.APIGatewayProxyRequest{
        Body: `{"username":"test","email":"test@example.com","password":"password"}`,
    }
    
    response, err := handler.RegisterUser(request)
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if response.StatusCode != 200 {
        t.Fatalf("Expected status 200, got %d", response.StatusCode)
    }
}
```

### Integration Testing
- Test with real AWS services
- Use local testing tools like SAM CLI
- Test with different event types

---

## Deployment Considerations

### 1. Build for Lambda
```bash
# Build for Linux (Lambda runtime)
GOOS=linux GOARCH=amd64 go build -o main main.go

# Create deployment package
zip deployment.zip main
```

### 2. Environment Variables
```go
// Use environment variables for configuration
databaseURL := os.Getenv("DATABASE_URL")
jwtSecret := os.Getenv("JWT_SECRET")
```

### 3. Cold Start Optimization
- Keep dependencies minimal
- Use connection pooling
- Initialize expensive resources once

---

## Key Takeaways

1. **Start with data structures** - Define your contracts first
2. **Build bottom-up** - Models → Utils → Middleware → Handlers → App → Main
3. **Use middleware pattern** - Separate concerns and make code composable
4. **Handle errors properly** - Always check and handle errors appropriately
5. **Test early and often** - Write tests as you build
6. **Keep it simple** - Don't over-engineer, start simple and iterate

---

## Next Steps for Your Project

1. **Implement the missing packages** using the step-by-step guide above
2. **Add proper error handling** to all functions
3. **Implement JWT authentication** in the middleware
4. **Add database integration** for user storage
5. **Write comprehensive tests** for all components
6. **Add logging and monitoring** for production use

Remember: **Start simple, test frequently, and iterate based on requirements!**
