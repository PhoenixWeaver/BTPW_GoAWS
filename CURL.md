
# 🧪 API Testing Commands

## Test Case 1: Register New User (Should Succeed)
```bash
curl -X POST https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/register \
  -H "Content-Type: application/json" \
  -d '{"username":"@ Firebird 1", "password":"@ Password123 !"}'
```
**Expected:** ✅ `"BINGO !!! Successfully, Registered Voter !"` (HTTP 200)

## Test Case 2: Login with Correct Credentials (Should Succeed)
```bash
curl -X POST https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login \
  -H "Content-Type: application/json" \
  -d '{"username":"@ Firebird 1", "password":"@ Password123 !"}'
```
**Expected:** ✅ `"Almmost successfull Login, JWT token generation is not implemented yet"` (HTTP 200)

## Test Case 3: Register Existing User (Should Fail)
```bash
curl -X POST https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/register \
  -H "Content-Type: application/json" \
  -d '{"username":"@ Firebird 1", "password":"Password123!"}'
```
**Expected:** ❌ `"Hey there. See again. User already exists"` (HTTP 409 Conflict)

## Test Case 4: Login with Wrong Password (Should Fail)
```bash
curl -X POST https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login \
  -H "Content-Type: application/json" \
  -d '{"username":"@ Firebird 1", "password":"Password123!"}'
```
**Expected:** ❌ `"Invalid login credentials"` (HTTP 401 Unauthorized)

## Test Case 5: Access Protected Route (JWT Not Implemented)
```bash
curl -X GET https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/protected \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN"
```
**Expected:** ❌ JWT validation error (JWT token generation not implemented yet)

---

# 🔧 PowerShell Testing Script

```powershell
# Complete API testing workflow
Write-Host "🔍 Testing API Endpoints..." -ForegroundColor Cyan

# Test Case 1: Register New User
Write-Host "1. Testing Registration..." -ForegroundColor Yellow
try {
    $reg = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/register" -Method POST -ContentType "application/json" -Body '{"username":"@ Firebird 1", "password":"@ Password123 !"}'
    Write-Host "✅ Register Success: $reg" -ForegroundColor Green
} catch {
    Write-Host "❌ Register Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Case 2: Login with Correct Credentials
Write-Host "2. Testing Login..." -ForegroundColor Yellow
try {
    $login = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login" -Method POST -ContentType "application/json" -Body '{"username":"@ Firebird 1", "password":"@ Password123 !"}'
    Write-Host "✅ Login Success: $login" -ForegroundColor Green
} catch {
    Write-Host "❌ Login Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Case 3: Try to Register Existing User (Should Fail)
Write-Host "3. Testing Duplicate Registration..." -ForegroundColor Yellow
try {
    $duplicate = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/register" -Method POST -ContentType "application/json" -Body '{"username":"@ Firebird 1", "password":"Password123!"}'
    Write-Host "❌ Unexpected Success: $duplicate" -ForegroundColor Red
} catch {
    Write-Host "✅ Duplicate Registration Correctly Blocked: $($_.Exception.Message)" -ForegroundColor Green
}

# Test Case 4: Login with Wrong Password (Should Fail)
Write-Host "4. Testing Wrong Password..." -ForegroundColor Yellow
try {
    $wrongPass = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login" -Method POST -ContentType "application/json" -Body '{"username":"@ Firebird 1", "password":"Password123!"}'
    Write-Host "❌ Unexpected Success: $wrongPass" -ForegroundColor Red
} catch {
    Write-Host "✅ Wrong Password Correctly Rejected: $($_.Exception.Message)" -ForegroundColor Green
}

# Test Case 5: Protected Route (JWT Not Implemented)
Write-Host "5. Testing Protected Route..." -ForegroundColor Yellow
try {
    $prot = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/protected" -Method GET
    Write-Host "❌ Unexpected Success: $prot" -ForegroundColor Red
} catch {
    Write-Host "✅ Protected Route Correctly Blocked: $($_.Exception.Message)" -ForegroundColor Green
}

Write-Host "🎉 Testing Complete!" -ForegroundColor Cyan

===========================================
## Test Case 6: Access Protected Route with Valid JWT Token (Should Succeed)
## Test Case 6: Access Protected Route with Valid JWT Token (Should Succeed)
```bash
# First, get a valid JWT token by logging in
curl -X POST https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login \
  -H "Content-Type: application/json" \
  -d '{"username":"@ Firebird 1", "password":"@ Password123 !"}' \
  -o login_response.json

# Extract the JWT token from the response (look for "access_token" field)
# Then use the token to access the protected route
curl -X GET https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/protected \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```
**Expected:** ✅ `"Protected content"` (HTTP 200) - **Currently getting Internal Server Error**

---

## **Updated PowerShell Testing Script**

```powershell
# Test Case 6: Access Protected Route with Valid JWT Token
Write-Host "6. Testing Protected Route with Valid Token..." -ForegroundColor Yellow

# Step 1: Login to get a valid JWT token
try {
    $loginResponse = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/login" -Method POST -ContentType "application/json" -Body '{"username":"@ Firebird 1", "password":"@ Password123 !"}'
    
    # Extract JWT token from response - FIXED: Use access_token instead of token
    $jwtToken = $loginResponse.access_token
    
    if ($jwtToken -and $jwtToken -ne "N/A") {
        Write-Host "✅ JWT Token Retrieved: $($jwtToken.Substring(0,20))..." -ForegroundColor Green
        
        # Step 2: Use the token to access protected route
        try {
            $headers = @{
                "Content-Type" = "application/json"
                "Authorization" = "Bearer $jwtToken"
            }
            
            $protectedResponse = Invoke-RestMethod -Uri "https://z8qo2a1k5f.execute-api.us-east-1.amazonaws.com/prod/protected" -Method GET -Headers $headers
            Write-Host "✅ Protected Route Access Success: $protectedResponse" -ForegroundColor Green
        } catch {
            Write-Host "❌ Protected Route Access Failed: $($_.Exception.Message)" -ForegroundColor Red
            Write-Host "🔍 This suggests the JWT middleware or protected route needs debugging" -ForegroundColor Yellow
        }
    } else {
        Write-Host "❌ No valid JWT token received from login" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Login failed, cannot test protected route: $($_.Exception.Message)" -ForegroundColor Red
}
```

## **🎯 Next Steps:**

The JWT token is being generated correctly, but there's an issue with the protected route implementation. You may need to:

1. **Check Lambda logs** in AWS CloudWatch
2. **Verify JWT middleware** is properly implemented
3. **Test the protected route** without JWT first to isolate the issue

The test case structure is now correct - it's just a matter of debugging the protected route implementation! 🚀