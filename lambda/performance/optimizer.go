package performance

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// PerformanceOptimizer provides performance optimization utilities
type PerformanceOptimizer struct {
	// Connection pooling
	dbClient *dynamodb.Client
	mu       sync.RWMutex
	
	// Caching
	cache    map[string]interface{}
	cacheMu  sync.RWMutex
	
	// Metrics
	requestCount int64
	avgLatency   time.Duration
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer() *PerformanceOptimizer {
	return &PerformanceOptimizer{
		cache: make(map[string]interface{}),
	}
}

// OptimizeDatabaseConnection optimizes database connections
func (po *PerformanceOptimizer) OptimizeDatabaseConnection() error {
	// Load AWS configuration with optimized settings
	cfg, err := config.LoadDefaultConfig(context.TODO(), 
		config.WithRetryMaxAttempts(3),
		config.WithRetryMode("adaptive"),
	)
	if err != nil {
		return err
	}

	// Create DynamoDB client with optimized settings
	po.dbClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		// Optimize connection settings
		o.RetryMaxAttempts = 3
		o.RetryMode = "adaptive"
	})

	return nil
}

// GetCachedValue retrieves a cached value
func (po *PerformanceOptimizer) GetCachedValue(key string) (interface{}, bool) {
	po.cacheMu.RLock()
	defer po.cacheMu.RUnlock()
	
	value, exists := po.cache[key]
	return value, exists
}

// SetCachedValue sets a cached value
func (po *PerformanceOptimizer) SetCachedValue(key string, value interface{}) {
	po.cacheMu.Lock()
	defer po.cacheMu.Unlock()
	
	po.cache[key] = value
}

// ClearCache clears the cache
func (po *PerformanceOptimizer) ClearCache() {
	po.cacheMu.Lock()
	defer po.cacheMu.Unlock()
	
	po.cache = make(map[string]interface{})
}

// GetDatabaseClient returns the optimized database client
func (po *PerformanceOptimizer) GetDatabaseClient() *dynamodb.Client {
	po.mu.RLock()
	defer po.mu.RUnlock()
	
	return po.dbClient
}

// RecordRequest records a request for metrics
func (po *PerformanceOptimizer) RecordRequest(latency time.Duration) {
	po.mu.Lock()
	defer po.mu.Unlock()
	
	po.requestCount++
	po.avgLatency = (po.avgLatency + latency) / 2
}

// GetMetrics returns performance metrics
func (po *PerformanceOptimizer) GetMetrics() map[string]interface{} {
	po.mu.RLock()
	defer po.mu.RUnlock()
	
	return map[string]interface{}{
		"request_count": po.requestCount,
		"avg_latency":   po.avgLatency.String(),
	}
}
