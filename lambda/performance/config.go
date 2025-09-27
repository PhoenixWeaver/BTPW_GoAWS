package performance

import (
	"os"
	"strconv"
	"time"
)

// PerformanceConfig holds performance configuration
type PerformanceConfig struct {
	// Lambda settings
	MemorySize    int
	Timeout       time.Duration
	Concurrency   int
	
	// Database settings
	MaxConnections    int
	ConnectionTimeout time.Duration
	QueryTimeout      time.Duration
	
	// Caching settings
	CacheEnabled      bool
	CacheTTL          time.Duration
	MaxCacheSize      int
	
	// Monitoring settings
	MetricsEnabled    bool
	LogLevel         string
	SamplingRate    float64
	
	// Optimization settings
	EnablePooling    bool
	EnableCompression bool
	EnableBatching   bool
}

// LoadPerformanceConfig loads performance configuration from environment variables
func LoadPerformanceConfig() *PerformanceConfig {
	config := &PerformanceConfig{
		// Default values
		MemorySize:         512,
		Timeout:            30 * time.Second,
		Concurrency:        10,
		MaxConnections:     10,
		ConnectionTimeout: 5 * time.Second,
		QueryTimeout:      10 * time.Second,
		CacheEnabled:       true,
		CacheTTL:           5 * time.Minute,
		MaxCacheSize:       1000,
		MetricsEnabled:     true,
		LogLevel:           "INFO",
		SamplingRate:       1.0,
		EnablePooling:      true,
		EnableCompression:  true,
		EnableBatching:     true,
	}
	
	// Load from environment variables
	if val := os.Getenv("LAMBDA_MEMORY_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.MemorySize = size
		}
	}
	
	if val := os.Getenv("LAMBDA_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.Timeout = timeout
		}
	}
	
	if val := os.Getenv("LAMBDA_CONCURRENCY"); val != "" {
		if concurrency, err := strconv.Atoi(val); err == nil {
			config.Concurrency = concurrency
		}
	}
	
	if val := os.Getenv("DB_MAX_CONNECTIONS"); val != "" {
		if maxConn, err := strconv.Atoi(val); err == nil {
			config.MaxConnections = maxConn
		}
	}
	
	if val := os.Getenv("DB_CONNECTION_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.ConnectionTimeout = timeout
		}
	}
	
	if val := os.Getenv("DB_QUERY_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.QueryTimeout = timeout
		}
	}
	
	if val := os.Getenv("CACHE_ENABLED"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.CacheEnabled = enabled
		}
	}
	
	if val := os.Getenv("CACHE_TTL"); val != "" {
		if ttl, err := time.ParseDuration(val); err == nil {
			config.CacheTTL = ttl
		}
	}
	
	if val := os.Getenv("CACHE_MAX_SIZE"); val != "" {
		if maxSize, err := strconv.Atoi(val); err == nil {
			config.MaxCacheSize = maxSize
		}
	}
	
	if val := os.Getenv("METRICS_ENABLED"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.MetricsEnabled = enabled
		}
	}
	
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.LogLevel = val
	}
	
	if val := os.Getenv("SAMPLING_RATE"); val != "" {
		if rate, err := strconv.ParseFloat(val, 64); err == nil {
			config.SamplingRate = rate
		}
	}
	
	if val := os.Getenv("ENABLE_POOLING"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnablePooling = enabled
		}
	}
	
	if val := os.Getenv("ENABLE_COMPRESSION"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableCompression = enabled
		}
	}
	
	if val := os.Getenv("ENABLE_BATCHING"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableBatching = enabled
		}
	}
	
	return config
}

// GetOptimalMemorySize calculates optimal memory size based on workload
func (pc *PerformanceConfig) GetOptimalMemorySize(workloadType string) int {
	switch workloadType {
	case "light":
		return 256
	case "medium":
		return 512
	case "heavy":
		return 1024
	case "intensive":
		return 2048
	default:
		return pc.MemorySize
	}
}

// GetOptimalTimeout calculates optimal timeout based on operation type
func (pc *PerformanceConfig) GetOptimalTimeout(operationType string) time.Duration {
	switch operationType {
	case "quick":
		return 5 * time.Second
	case "standard":
		return 15 * time.Second
	case "complex":
		return 30 * time.Second
	case "batch":
		return 60 * time.Second
	default:
		return pc.Timeout
	}
}

// IsOptimizationEnabled checks if a specific optimization is enabled
func (pc *PerformanceConfig) IsOptimizationEnabled(optimization string) bool {
	switch optimization {
	case "pooling":
		return pc.EnablePooling
	case "compression":
		return pc.EnableCompression
	case "batching":
		return pc.EnableBatching
	case "caching":
		return pc.CacheEnabled
	case "metrics":
		return pc.MetricsEnabled
	default:
		return false
	}
}
