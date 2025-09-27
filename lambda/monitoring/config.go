package monitoring

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	// CloudWatch settings
	LogGroup      string
	LogStream     string
	Namespace     string
	DashboardName string

	// Logging settings
	LogLevel         string
	EnableCloudWatch bool
	EnableConsole    bool

	// Metrics settings
	EnableMetrics   bool
	MetricsInterval time.Duration
	CustomMetrics   bool

	// Health check settings
	EnableHealthCheck   bool
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration

	// Dashboard settings
	EnableDashboard     bool
	AutoCreateDashboard bool

	// Alarms settings
	EnableAlarms    bool
	AlarmActions    []string
	AlarmThresholds map[string]float64

	// Performance settings
	EnablePerformanceMonitoring bool
	PerformanceThresholds       map[string]float64
}

// LoadMonitoringConfig loads monitoring configuration from environment variables
func LoadMonitoringConfig() *MonitoringConfig {
	config := &MonitoringConfig{
		// Default values
		LogGroup:                    "/aws/lambda/function",
		LogStream:                   "lambda-function",
		Namespace:                   "AWS/Lambda",
		DashboardName:               "LambdaFunctionDashboard",
		LogLevel:                    "INFO",
		EnableCloudWatch:            true,
		EnableConsole:               true,
		EnableMetrics:               true,
		MetricsInterval:             5 * time.Minute,
		CustomMetrics:               true,
		EnableHealthCheck:           true,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckTimeout:          10 * time.Second,
		EnableDashboard:             true,
		AutoCreateDashboard:         true,
		EnableAlarms:                true,
		AlarmActions:                []string{},
		AlarmThresholds:             make(map[string]float64),
		EnablePerformanceMonitoring: true,
		PerformanceThresholds:       make(map[string]float64),
	}

	// Load from environment variables
	if val := os.Getenv("LOG_GROUP"); val != "" {
		config.LogGroup = val
	}

	if val := os.Getenv("LOG_STREAM"); val != "" {
		config.LogStream = val
	}

	if val := os.Getenv("METRICS_NAMESPACE"); val != "" {
		config.Namespace = val
	}

	if val := os.Getenv("DASHBOARD_NAME"); val != "" {
		config.DashboardName = val
	}

	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.LogLevel = val
	}

	if val := os.Getenv("ENABLE_CLOUDWATCH"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableCloudWatch = enabled
		}
	}

	if val := os.Getenv("ENABLE_CONSOLE"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableConsole = enabled
		}
	}

	if val := os.Getenv("ENABLE_METRICS"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableMetrics = enabled
		}
	}

	if val := os.Getenv("METRICS_INTERVAL"); val != "" {
		if interval, err := time.ParseDuration(val); err == nil {
			config.MetricsInterval = interval
		}
	}

	if val := os.Getenv("CUSTOM_METRICS"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.CustomMetrics = enabled
		}
	}

	if val := os.Getenv("ENABLE_HEALTH_CHECK"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableHealthCheck = enabled
		}
	}

	if val := os.Getenv("HEALTH_CHECK_INTERVAL"); val != "" {
		if interval, err := time.ParseDuration(val); err == nil {
			config.HealthCheckInterval = interval
		}
	}

	if val := os.Getenv("HEALTH_CHECK_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.HealthCheckTimeout = timeout
		}
	}

	if val := os.Getenv("ENABLE_DASHBOARD"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableDashboard = enabled
		}
	}

	if val := os.Getenv("AUTO_CREATE_DASHBOARD"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.AutoCreateDashboard = enabled
		}
	}

	if val := os.Getenv("ENABLE_ALARMS"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnableAlarms = enabled
		}
	}

	if val := os.Getenv("ENABLE_PERFORMANCE_MONITORING"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.EnablePerformanceMonitoring = enabled
		}
	}

	// Set default alarm thresholds
	config.AlarmThresholds = map[string]float64{
		"error_rate":      10.0,
		"response_time":   5.0,
		"memory_usage":    100000000, // 100MB
		"auth_failures":   5.0,
		"security_events": 1.0,
	}

	// Set default performance thresholds
	config.PerformanceThresholds = map[string]float64{
		"max_response_time": 5.0,
		"max_memory_usage":  100000000, // 100MB
		"max_cpu_usage":     80.0,
		"max_error_rate":    5.0,
	}

	return config
}

// MonitoringManager provides centralized monitoring management
type MonitoringManager struct {
	config           *MonitoringConfig
	logger           *Logger
	metricsCollector *MetricsCollector
	healthChecker    *HealthChecker
	dashboardConfig  *DashboardConfig
	ctx              context.Context
}

// NewMonitoringManager creates a new monitoring manager
func NewMonitoringManager(ctx context.Context) *MonitoringManager {
	monitoringConfig := LoadMonitoringConfig()

	// Create AWS clients
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	cloudwatchClient := cloudwatch.NewFromConfig(awsConfig)
	cloudwatchLogsClient := cloudwatchlogs.NewFromConfig(awsConfig)

	// Create components
	logger := NewLogger(ctx, cloudwatchLogsClient, monitoringConfig.LogGroup, monitoringConfig.LogStream)
	metricsCollector := NewMetricsCollector(cloudwatchClient, monitoringConfig.Namespace)
	healthChecker := NewHealthChecker(ctx)
	dashboardConfig := NewDashboardConfig(cloudwatchClient, monitoringConfig.DashboardName, monitoringConfig.Namespace)

	return &MonitoringManager{
		config:           monitoringConfig,
		logger:           logger,
		metricsCollector: metricsCollector,
		healthChecker:    healthChecker,
		dashboardConfig:  dashboardConfig,
		ctx:              ctx,
	}
}

// Initialize initializes the monitoring system
func (mm *MonitoringManager) Initialize() error {
	// Set log level
	switch mm.config.LogLevel {
	case "DEBUG":
		mm.logger.SetLevel(DEBUG)
	case "INFO":
		mm.logger.SetLevel(INFO)
	case "WARN":
		mm.logger.SetLevel(WARN)
	case "ERROR":
		mm.logger.SetLevel(ERROR)
	case "FATAL":
		mm.logger.SetLevel(FATAL)
	}

	// Create dashboard if enabled
	if mm.config.EnableDashboard && mm.config.AutoCreateDashboard {
		if err := mm.dashboardConfig.CreateDashboard(mm.ctx); err != nil {
			mm.logger.Error("Failed to create dashboard", err)
		}
	}

	// Create alarms if enabled
	if mm.config.EnableAlarms {
		if err := mm.dashboardConfig.CreateAlarms(mm.ctx); err != nil {
			mm.logger.Error("Failed to create alarms", err)
		}
	}

	mm.logger.Info("Monitoring system initialized")
	return nil
}

// GetLogger returns the logger
func (mm *MonitoringManager) GetLogger() *Logger {
	return mm.logger
}

// GetMetricsCollector returns the metrics collector
func (mm *MonitoringManager) GetMetricsCollector() *MetricsCollector {
	return mm.metricsCollector
}

// GetHealthChecker returns the health checker
func (mm *MonitoringManager) GetHealthChecker() *HealthChecker {
	return mm.healthChecker
}

// GetDashboardConfig returns the dashboard config
func (mm *MonitoringManager) GetDashboardConfig() *DashboardConfig {
	return mm.dashboardConfig
}

// GetConfig returns the monitoring config
func (mm *MonitoringManager) GetConfig() *MonitoringConfig {
	return mm.config
}

// Shutdown shuts down the monitoring system
func (mm *MonitoringManager) Shutdown() {
	mm.logger.Info("Monitoring system shutting down")
}
