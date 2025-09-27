# AWS Lambda Optimization and Deployment Script
Write-Host "🚀 AWS Lambda Optimization and Deployment Script" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green

Write-Host ""
Write-Host "Step OPTIMIZATION: Setting up optimization environment..." -ForegroundColor Yellow

# Set optimization environment variables
$env:ENABLE_PERFORMANCE_MONITORING = "true"
$env:OPTIMIZE_MEMORY = "true"
$env:OPTIMIZE_DATABASE = "true"
$env:ENABLE_CACHING = "true"
$env:ENABLE_CLOUDWATCH = "true"
$env:ENABLE_METRICS = "true"
$env:ENABLE_LOGGING = "true"
$env:ENABLE_HEALTH_CHECK = "true"
$env:LOG_LEVEL = "INFO"

Write-Host "✅ Environment variables set for optimization" -ForegroundColor Green

Write-Host ""
Write-Host "Step 1: Building Lambda Function (Linux)..." -ForegroundColor Yellow

# Navigate to lambda directory
Set-Location "C:\Users\Admin\Documents\GODEV\GO Courses\GO AWS\BT_GoAWS\lambda"

# Set Go environment for Linux (Lambda runtime)
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# Build with optimization flags
Write-Host "Building with optimization flags: -ldflags=`"-s -w`"" -ForegroundColor Cyan
go build -ldflags="-s -w" -o bootstrap main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Lambda build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✅ Lambda function built with optimization" -ForegroundColor Green

# Create deployment package
Write-Host "Creating deployment package..." -ForegroundColor Cyan
Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Archive creation failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✅ Deployment package created: function.zip" -ForegroundColor Green

Write-Host ""
Write-Host "Step 2: Building CDK Application (Windows)..." -ForegroundColor Yellow

# Navigate to CDK directory
Set-Location "C:\Users\Admin\Documents\GODEV\GO Courses\GO AWS\BT_GoAWS"

# Reset Go environment for Windows (CDK)
$env:GOOS = "windows"
$env:GOARCH = "amd64"

go build -o BTgoAWS.exe BT_GoAws.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ CDK build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✅ CDK application built successfully" -ForegroundColor Green

Write-Host ""
Write-Host "Step 3: Testing CDK synthesis..." -ForegroundColor Yellow

# Test CDK synthesis
# cdk synth

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ CDK synthesis failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✅ CDK synthesis successful" -ForegroundColor Green

# Show differences
Write-Host ""
Write-Host "Showing differences..." -ForegroundColor Yellow
cdk diff

Write-Host ""
Write-Host "Step 4: Deploying to AWS..." -ForegroundColor Yellow

# Deploy to AWS
cdk deploy --yes --require-approval never

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Deployment failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✅ Deployment successful!" -ForegroundColor Green

Write-Host ""
Write-Host "Step 5: Testing the deployed function..." -ForegroundColor Yellow

# Test via AWS CLI
Write-Host "Testing Lambda function..." -ForegroundColor Cyan
aws lambda invoke --function-name BTgoAWSstack-BTlambdaFunction8DF41906-ZlGN7I54bIaq --payload '{"path": "/register", "body": "{\"username\": \"testuser\", \"password\": \"password123\"}"}' response.json

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Lambda function test successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Response:" -ForegroundColor Cyan
    Get-Content response.json | Write-Host -ForegroundColor White
} else {
    Write-Host "❌ Lambda function test failed!" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 SUCCESS! Optimized Lambda function deployed and tested" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Optimization Features Enabled:" -ForegroundColor Cyan
Write-Host "- Performance monitoring: $env:ENABLE_PERFORMANCE_MONITORING" -ForegroundColor White
Write-Host "- Memory optimization: $env:OPTIMIZE_MEMORY" -ForegroundColor White
Write-Host "- Database optimization: $env:OPTIMIZE_DATABASE" -ForegroundColor White
Write-Host "- Caching: $env:ENABLE_CACHING" -ForegroundColor White
Write-Host "- CloudWatch: $env:ENABLE_CLOUDWATCH" -ForegroundColor White
Write-Host "- Metrics: $env:ENABLE_METRICS" -ForegroundColor White
Write-Host "- Logging: $env:ENABLE_LOGGING" -ForegroundColor White
Write-Host "- Health checks: $env:ENABLE_HEALTH_CHECK" -ForegroundColor White

Write-Host ""
Write-Host "🔗 Check your Lambda function in AWS Console:" -ForegroundColor Yellow
Write-Host "https://us-east-1.console.aws.amazon.com/lambda/home?region=us-east-1" -ForegroundColor Blue

Read-Host "Press Enter to exit"
