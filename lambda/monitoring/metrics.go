package monitoring

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// MetricsCollector provides metrics collection capabilities
type MetricsCollector struct {
	cloudwatchClient *cloudwatch.Client
	namespace        string
	metrics          map[string]float64
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(cloudwatchClient *cloudwatch.Client, namespace string) *MetricsCollector {
	return &MetricsCollector{
		cloudwatchClient: cloudwatchClient,
		namespace:        namespace,
		metrics:          make(map[string]float64),
	}
}

// MetricType represents the type of metric
type MetricType string

const (
	Counter   MetricType = "Counter"
	Gauge     MetricType = "Gauge"
	Histogram MetricType = "Histogram"
	Summary   MetricType = "Summary"
)

// Metric represents a metric
type Metric struct {
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
	Tags      map[string]string
}

// RecordMetric records a metric
func (mc *MetricsCollector) RecordMetric(name string, value float64, unit string) {
	mc.metrics[name] = value
	
	// Create metric datum
	metric := Metric{
		Name:      name,
		Value:     value,
		Unit:      unit,
		Timestamp: time.Now(),
		Tags:      make(map[string]string),
	}
	
	// Send to CloudWatch
	mc.sendMetricToCloudWatch(metric)
}

// RecordCounter records a counter metric
func (mc *MetricsCollector) RecordCounter(name string, value float64) {
	mc.RecordMetric(name, value, "Count")
}

// RecordGauge records a gauge metric
func (mc *MetricsCollector) RecordGauge(name string, value float64) {
	mc.RecordMetric(name, value, "None")
}

// RecordHistogram records a histogram metric
func (mc *MetricsCollector) RecordHistogram(name string, value float64) {
	mc.RecordMetric(name, value, "None")
}

// RecordSummary records a summary metric
func (mc *MetricsCollector) RecordSummary(name string, value float64) {
	mc.RecordMetric(name, value, "None")
}

// RecordRequestMetrics records request-related metrics
func (mc *MetricsCollector) RecordRequestMetrics(method, path string, statusCode int, duration time.Duration) {
	// Record request count
	mc.RecordCounter("requests_total", 1)
	
	// Record request duration
	mc.RecordHistogram("request_duration_seconds", duration.Seconds())
	
	// Record status code
	mc.RecordCounter("requests_by_status", float64(statusCode))
	
	// Record by method
	mc.RecordCounter("requests_by_method", 1)
	
	// Record by path
	mc.RecordCounter("requests_by_path", 1)
}

// RecordDatabaseMetrics records database-related metrics
func (mc *MetricsCollector) RecordDatabaseMetrics(operation, table string, duration time.Duration, success bool) {
	// Record database operation count
	mc.RecordCounter("database_operations_total", 1)
	
	// Record database operation duration
	mc.RecordHistogram("database_operation_duration_seconds", duration.Seconds())
	
	// Record success/failure
	if success {
		mc.RecordCounter("database_operations_success", 1)
	} else {
		mc.RecordCounter("database_operations_failure", 1)
	}
	
	// Record by operation type
	mc.RecordCounter("database_operations_by_type", 1)
	
	// Record by table
	mc.RecordCounter("database_operations_by_table", 1)
}

// RecordAuthenticationMetrics records authentication-related metrics
func (mc *MetricsCollector) RecordAuthenticationMetrics(username string, success bool, duration time.Duration) {
	// Record authentication attempts
	mc.RecordCounter("authentication_attempts_total", 1)
	
	// Record authentication duration
	mc.RecordHistogram("authentication_duration_seconds", duration.Seconds())
	
	// Record success/failure
	if success {
		mc.RecordCounter("authentication_success", 1)
	} else {
		mc.RecordCounter("authentication_failure", 1)
	}
}

// RecordSecurityMetrics records security-related metrics
func (mc *MetricsCollector) RecordSecurityMetrics(event string, severity string) {
	// Record security events
	mc.RecordCounter("security_events_total", 1)
	
	// Record by event type
	mc.RecordCounter("security_events_by_type", 1)
	
	// Record by severity
	mc.RecordCounter("security_events_by_severity", 1)
}

// RecordPerformanceMetrics records performance-related metrics
func (mc *MetricsCollector) RecordPerformanceMetrics(operation string, duration time.Duration, memoryUsage int64) {
	// Record operation duration
	mc.RecordHistogram("operation_duration_seconds", duration.Seconds())
	
	// Record memory usage
	mc.RecordGauge("memory_usage_bytes", float64(memoryUsage))
	
	// Record by operation
	mc.RecordCounter("operations_by_type", 1)
}

// RecordErrorMetrics records error-related metrics
func (mc *MetricsCollector) RecordErrorMetrics(errorType, errorMessage string) {
	// Record error count
	mc.RecordCounter("errors_total", 1)
	
	// Record by error type
	mc.RecordCounter("errors_by_type", 1)
	
	// Record by error message
	mc.RecordCounter("errors_by_message", 1)
}

// sendMetricToCloudWatch sends a metric to CloudWatch
func (mc *MetricsCollector) sendMetricToCloudWatch(metric Metric) {
	if mc.cloudwatchClient == nil {
		return
	}
	
	// Create metric datum
	metricDatum := types.MetricDatum{
		MetricName: &metric.Name,
		Value:      &metric.Value,
		Unit:       types.StandardUnit(metric.Unit),
		Timestamp:  &metric.Timestamp,
	}
	
	// Add dimensions
	if len(metric.Tags) > 0 {
		var dimensions []types.Dimension
		for key, value := range metric.Tags {
			dimensions = append(dimensions, types.Dimension{
				Name:  &key,
				Value: &value,
			})
		}
		metricDatum.Dimensions = dimensions
	}
	
	// Send to CloudWatch
	_, err := mc.cloudwatchClient.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace:  &mc.namespace,
		MetricData: []types.MetricDatum{metricDatum},
	})
	
	if err != nil {
		// Log error but don't fail
		// In production, you might want to use a fallback mechanism
	}
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() map[string]float64 {
	return mc.metrics
}

// ClearMetrics clears current metrics
func (mc *MetricsCollector) ClearMetrics() {
	mc.metrics = make(map[string]float64)
}

// CustomMetricsCollector provides custom metrics collection
type CustomMetricsCollector struct {
	*MetricsCollector
	customMetrics map[string]interface{}
}

// NewCustomMetricsCollector creates a new custom metrics collector
func NewCustomMetricsCollector(cloudwatchClient *cloudwatch.Client, namespace string) *CustomMetricsCollector {
	return &CustomMetricsCollector{
		MetricsCollector: NewMetricsCollector(cloudwatchClient, namespace),
		customMetrics:    make(map[string]interface{}),
	}
}

// RecordCustomMetric records a custom metric
func (cmc *CustomMetricsCollector) RecordCustomMetric(name string, value interface{}) {
	cmc.customMetrics[name] = value
	
	// Convert to float64 if possible
	if floatValue, ok := value.(float64); ok {
		cmc.RecordMetric(name, floatValue, "None")
	}
}

// GetCustomMetrics returns custom metrics
func (cmc *CustomMetricsCollector) GetCustomMetrics() map[string]interface{} {
	return cmc.customMetrics
}

// ClearCustomMetrics clears custom metrics
func (cmc *CustomMetricsCollector) ClearCustomMetrics() {
	cmc.customMetrics = make(map[string]interface{})
}
