# 🚀 How to Run Optimization and Enable Monitoring

## 📋 Quick Start Guide

### 1. **Running Performance Optimization**

#### Step 1: Navigate to Lambda Directory
```bash
cd C:\Users\Admin\Documents\GODEV\GO Courses\GO AWS\BT_GoAWS\lambda
```

#### Step 2: Enable Performance Optimization
```bash
# Set environment variables for optimization
export ENABLE_PERFORMANCE_MONITORING=true
export OPTIMIZE_MEMORY=true
export OPTIMIZE_DATABASE=true
export ENABLE_CACHING=true

# Or on Windows PowerShell:
$env:ENABLE_PERFORMANCE_MONITORING="true"
$env:OPTIMIZE_MEMORY="true"
$env:OPTIMIZE_DATABASE="true"
$env:ENABLE_CACHING="true"
```

#### Step 3: Run with Optimization
```bash
# Build with optimization flags
go build -ldflags="-s -w" -o main main.go

# Run with performance monitoring
./main
```

### 2. **Enabling Monitoring System**

#### Step 1: Set Monitoring Environment Variables
```bash
# Enable all monitoring features
export ENABLE_CLOUDWATCH=true
export ENABLE_METRICS=true
export ENABLE_LOGGING=true
export ENABLE_HEALTH_CHECK=true
export ENABLE_DASHBOARD=true

# Logging configuration
export LOG_LEVEL=INFO
export LOG_GROUP="/aws/lambda/function"
export LOG_STREAM="lambda-function"

# Metrics configuration
export METRICS_NAMESPACE="AWS/Lambda"
export ENABLE_CUSTOM_METRICS=true

# Health check configuration
export ENABLE_HEALTH_CHECK=true
export HEALTH_CHECK_INTERVAL=30s
export HEALTH_CHECK_TIMEOUT=10s
```

#### Step 2: Initialize Monitoring
```bash
# Run the application with monitoring enabled
go run main.go
```

### 3. **Using the Makefile Commands**

#### Available Commands
```bash
# Build with optimization
make build

# Run tests with monitoring
make test

# Deploy with monitoring enabled
make deploy

# Run performance tests
make perf-test

# Start monitoring dashboard
make monitor

# Clean build artifacts
make clean
```

### 4. **Manual Optimization Commands**

#### Memory Optimization
```bash
# Run memory optimization
go run performance/memory_optimizer.go

# Monitor memory usage
go run performance/monitor.go --memory
```

#### Database Optimization
```bash
# Run database optimization
go run performance/database_optimizer.go

# Monitor database performance
go run performance/monitor.go --database
```

#### Performance Monitoring
```bash
# Start performance monitoring
go run performance/monitor.go --start

# View performance metrics
go run performance/monitor.go --metrics
```

### 5. **Monitoring Dashboard Setup**

#### Step 1: Create CloudWatch Dashboard
```bash
# Run dashboard configuration
go run monitoring/dashboard.go

# Or use AWS CLI
aws cloudwatch put-dashboard --dashboard-name "LambdaFunctionDashboard" --dashboard-body file://dashboard.json
```

#### Step 2: View Monitoring Data
```bash
# Check logs
aws logs describe-log-groups --log-group-name-prefix "/aws/lambda"

# View metrics
aws cloudwatch get-metric-statistics --namespace "AWS/Lambda" --metric-name "Duration"
```

### 6. **Health Check Endpoints**

#### Available Health Endpoints
```bash
# Liveness check
curl http://localhost:8080/health/live

# Readiness check
curl http://localhost:8080/health/ready

# Full health check
curl http://localhost:8080/health
```

### 7. **Real-time Monitoring**

#### Start Monitoring Session
```bash
# Start comprehensive monitoring
go run monitoring/config.go

# Monitor specific components
go run monitoring/logger.go --level INFO
go run monitoring/metrics.go --interval 30s
go run monitoring/health.go --timeout 10s
```

### 8. **Performance Testing**

#### Load Testing
```bash
# Run performance tests
go test -v ./performance -run TestPerformance

# Run load tests
go test -v ./e2e_test.go -run TestE2E_Performance
```

#### Benchmark Testing
```bash
# Run benchmarks
go test -bench=. ./performance

# Memory benchmarks
go test -bench=BenchmarkMemory ./performance
```

### 9. **Monitoring Configuration**

#### Create Monitoring Config File
```bash
# Create config file
cat > monitoring_config.json << EOF
{
  "log_level": "INFO",
  "enable_cloudwatch": true,
  "enable_metrics": true,
  "enable_health_check": true,
  "metrics_interval": "30s",
  "health_check_interval": "30s"
}
EOF
```

#### Load Configuration
```bash
# Load configuration
export MONITORING_CONFIG_FILE="monitoring_config.json"
go run main.go
```

### 10. **Troubleshooting**

#### Check Monitoring Status
```bash
# Check if monitoring is enabled
echo $ENABLE_CLOUDWATCH
echo $ENABLE_METRICS
echo $ENABLE_HEALTH_CHECK

# Check logs
tail -f /var/log/lambda.log
```

#### Common Issues
```bash
# If monitoring fails, check AWS credentials
aws sts get-caller-identity

# If performance optimization fails, check memory
go run performance/monitor.go --check-memory

# If health checks fail, check database connection
go run monitoring/health.go --check-database
```

## 🎯 **Quick Commands Summary**

### Enable Everything at Once
```bash
# Set all environment variables
export ENABLE_PERFORMANCE_MONITORING=true
export ENABLE_CLOUDWATCH=true
export ENABLE_METRICS=true
export ENABLE_LOGGING=true
export ENABLE_HEALTH_CHECK=true
export ENABLE_DASHBOARD=true

# Run with all features
go run main.go
```

### Check Status
```bash
# Check if everything is running
curl http://localhost:8080/health
curl http://localhost:8080/metrics
curl http://localhost:8080/status
```

### Stop Monitoring
```bash
# Stop all monitoring
unset ENABLE_CLOUDWATCH
unset ENABLE_METRICS
unset ENABLE_HEALTH_CHECK
```

## 📊 **Monitoring Output Examples**

### Performance Metrics
```
Memory Usage: 45MB
Response Time: 120ms
Database Latency: 25ms
Throughput: 150 req/min
```

### Health Status
```
Status: healthy
Database: connected
Memory: normal
CPU: 15%
```

### Log Output
```
[INFO] Performance monitoring enabled
[INFO] CloudWatch logging enabled
[INFO] Health checks running
[INFO] Metrics collection active
```

This guide will help you run optimization and enable monitoring for your AWS Lambda project! 🚀
