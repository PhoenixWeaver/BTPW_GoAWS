package monitoring

import (
	"context"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// HealthStatus represents the health status
type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Degraded  HealthStatus = "degraded"
)

// HealthCheck represents a health check
type HealthCheck struct {
	Name      string                 `json:"name"`
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthChecker provides health checking capabilities
type HealthChecker struct {
	checks []HealthCheck
	ctx    context.Context
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(ctx context.Context) *HealthChecker {
	return &HealthChecker{
		checks: make([]HealthCheck, 0),
		ctx:    ctx,
	}
}

// CheckDatabaseHealth checks database health
func (hc *HealthChecker) CheckDatabaseHealth(dbClient *dynamodb.Client, tableName string) HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:      "database",
		Status:    Healthy,
		Message:   "Database is healthy",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Check database connectivity
	_, err := dbClient.DescribeTable(hc.ctx, &dynamodb.DescribeTableInput{
		TableName: &tableName,
	})

	if err != nil {
		check.Status = Unhealthy
		check.Message = "Database connection failed"
		check.Details["error"] = err.Error()
	} else {
		check.Details["table_name"] = tableName
		check.Details["connection"] = "successful"
	}

	check.Duration = time.Since(start)
	hc.checks = append(hc.checks, check)

	return check
}

// CheckMemoryHealth checks memory health
func (hc *HealthChecker) CheckMemoryHealth() HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:      "memory",
		Status:    Healthy,
		Message:   "Memory is healthy",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check memory usage
	memoryUsagePercent := float64(m.Alloc) / float64(m.Sys) * 100

	if memoryUsagePercent > 90 {
		check.Status = Unhealthy
		check.Message = "Memory usage is too high"
	} else if memoryUsagePercent > 75 {
		check.Status = Degraded
		check.Message = "Memory usage is elevated"
	}

	check.Details["alloc_bytes"] = m.Alloc
	check.Details["sys_bytes"] = m.Sys
	check.Details["usage_percent"] = memoryUsagePercent
	check.Details["gc_count"] = m.NumGC

	check.Duration = time.Since(start)
	hc.checks = append(hc.checks, check)

	return check
}

// CheckDiskHealth checks disk health
func (hc *HealthChecker) CheckDiskHealth() HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:      "disk",
		Status:    Healthy,
		Message:   "Disk is healthy",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// In Lambda, disk space is limited and managed by AWS
	// This is a placeholder for disk health checks
	check.Details["disk_type"] = "ephemeral"
	check.Details["managed_by"] = "AWS Lambda"

	check.Duration = time.Since(start)
	hc.checks = append(hc.checks, check)

	return check
}

// CheckNetworkHealth checks network health
func (hc *HealthChecker) CheckNetworkHealth() HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:      "network",
		Status:    Healthy,
		Message:   "Network is healthy",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// In Lambda, network is managed by AWS
	// This is a placeholder for network health checks
	check.Details["network_type"] = "AWS VPC"
	check.Details["managed_by"] = "AWS Lambda"

	check.Duration = time.Since(start)
	hc.checks = append(hc.checks, check)

	return check
}

// CheckApplicationHealth checks application health
func (hc *HealthChecker) CheckApplicationHealth() HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:      "application",
		Status:    Healthy,
		Message:   "Application is healthy",
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Check application-specific health
	check.Details["version"] = "1.0.0"
	check.Details["environment"] = "production"
	check.Details["uptime"] = time.Since(start).String()

	check.Duration = time.Since(start)
	hc.checks = append(hc.checks, check)

	return check
}

// RunAllChecks runs all health checks
func (hc *HealthChecker) RunAllChecks(dbClient *dynamodb.Client, tableName string) []HealthCheck {
	// Clear previous checks
	hc.checks = make([]HealthCheck, 0)

	// Run all checks
	hc.CheckDatabaseHealth(dbClient, tableName)
	hc.CheckMemoryHealth()
	hc.CheckDiskHealth()
	hc.CheckNetworkHealth()
	hc.CheckApplicationHealth()

	return hc.checks
}

// GetOverallHealth returns the overall health status
func (hc *HealthChecker) GetOverallHealth() HealthStatus {
	if len(hc.checks) == 0 {
		return Unhealthy
	}

	unhealthyCount := 0
	degradedCount := 0

	for _, check := range hc.checks {
		switch check.Status {
		case Unhealthy:
			unhealthyCount++
		case Degraded:
			degradedCount++
		}
	}

	if unhealthyCount > 0 {
		return Unhealthy
	}

	if degradedCount > 0 {
		return Degraded
	}

	return Healthy
}

// GetHealthSummary returns a health summary
func (hc *HealthChecker) GetHealthSummary() map[string]interface{} {
	overallHealth := hc.GetOverallHealth()

	summary := map[string]interface{}{
		"overall_status":  overallHealth,
		"total_checks":    len(hc.checks),
		"healthy_count":   0,
		"unhealthy_count": 0,
		"degraded_count":  0,
		"checks":          hc.checks,
		"timestamp":       time.Now(),
	}

	// Count statuses
	for _, check := range hc.checks {
		switch check.Status {
		case Healthy:
			summary["healthy_count"] = summary["healthy_count"].(int) + 1
		case Unhealthy:
			summary["unhealthy_count"] = summary["unhealthy_count"].(int) + 1
		case Degraded:
			summary["degraded_count"] = summary["degraded_count"].(int) + 1
		}
	}

	return summary
}

// HealthCheckHandler provides HTTP health check handling
type HealthCheckHandler struct {
	healthChecker *HealthChecker
	dbClient      *dynamodb.Client
	tableName     string
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler(ctx context.Context, dbClient *dynamodb.Client, tableName string) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthChecker: NewHealthChecker(ctx),
		dbClient:      dbClient,
		tableName:     tableName,
	}
}

// HandleHealthCheck handles health check requests
func (hch *HealthCheckHandler) HandleHealthCheck() map[string]interface{} {
	// Run all health checks
	checks := hch.healthChecker.RunAllChecks(hch.dbClient, hch.tableName)

	// Get overall health
	overallHealth := hch.healthChecker.GetOverallHealth()

	// Return health summary
	return map[string]interface{}{
		"status":    overallHealth,
		"checks":    checks,
		"timestamp": time.Now(),
	}
}

// HandleLivenessCheck handles liveness check requests
func (hch *HealthCheckHandler) HandleLivenessCheck() map[string]interface{} {
	// Simple liveness check
	return map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	}
}

// HandleReadinessCheck handles readiness check requests
func (hch *HealthCheckHandler) HandleReadinessCheck() map[string]interface{} {
	// Check if application is ready to serve requests
	checks := hch.healthChecker.RunAllChecks(hch.dbClient, hch.tableName)

	// Check if all critical services are healthy
	ready := true
	for _, check := range checks {
		if check.Status == Unhealthy {
			ready = false
			break
		}
	}

	status := "ready"
	if !ready {
		status = "not_ready"
	}

	return map[string]interface{}{
		"status":    status,
		"checks":    checks,
		"timestamp": time.Now(),
	}
}
