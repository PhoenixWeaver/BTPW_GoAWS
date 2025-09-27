package performance

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// PerformanceMonitor provides performance monitoring capabilities
type PerformanceMonitor struct {
	cloudwatchClient *cloudwatch.Client
	metrics          map[string]float64
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(client *cloudwatch.Client) *PerformanceMonitor {
	return &PerformanceMonitor{
		cloudwatchClient: client,
		metrics:          make(map[string]float64),
	}
}

// RecordLatency records request latency
func (pm *PerformanceMonitor) RecordLatency(operation string, latency time.Duration) {
	pm.metrics[operation+"_latency"] = float64(latency.Milliseconds())
}

// RecordThroughput records request throughput
func (pm *PerformanceMonitor) RecordThroughput(operation string, count int) {
	pm.metrics[operation+"_throughput"] = float64(count)
}

// RecordError records error count
func (pm *PerformanceMonitor) RecordError(operation string, errorCount int) {
	pm.metrics[operation+"_errors"] = float64(errorCount)
}

// RecordMemoryUsage records memory usage
func (pm *PerformanceMonitor) RecordMemoryUsage(usage int64) {
	pm.metrics["memory_usage"] = float64(usage)
}

// SendMetricsToCloudWatch sends metrics to CloudWatch
func (pm *PerformanceMonitor) SendMetricsToCloudWatch(ctx context.Context, namespace string) error {
	if pm.cloudwatchClient == nil {
		return nil // Skip if no CloudWatch client
	}

	var metricData []types.MetricDatum
	timestamp := time.Now()

	for metricName, value := range pm.metrics {
		metricData = append(metricData, types.MetricDatum{
			MetricName: &metricName,
			Value:      &value,
			Timestamp:  &timestamp,
			Unit:       types.StandardUnitCount,
		})
	}

	if len(metricData) == 0 {
		return nil
	}

	_, err := pm.cloudwatchClient.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		Namespace:  &namespace,
		MetricData: metricData,
	})

	return err
}

// GetMetrics returns current metrics
func (pm *PerformanceMonitor) GetMetrics() map[string]float64 {
	return pm.metrics
}

// ClearMetrics clears current metrics
func (pm *PerformanceMonitor) ClearMetrics() {
	pm.metrics = make(map[string]float64)
}

// RecordCustomMetric records a custom metric
func (pm *PerformanceMonitor) RecordCustomMetric(metricName string, value float64) {
	pm.metrics[metricName] = value
}

// GetMetricValue gets a specific metric value
func (pm *PerformanceMonitor) GetMetricValue(metricName string) (float64, bool) {
	value, exists := pm.metrics[metricName]
	return value, exists
}

// RecordDatabaseOperation records database operation metrics
func (pm *PerformanceMonitor) RecordDatabaseOperation(operation string, latency time.Duration, success bool) {
	pm.RecordLatency("db_"+operation, latency)
	
	if success {
		pm.RecordThroughput("db_"+operation, 1)
	} else {
		pm.RecordError("db_"+operation, 1)
	}
}

// RecordAPIOperation records API operation metrics
func (pm *PerformanceMonitor) RecordAPIOperation(operation string, latency time.Duration, success bool) {
	pm.RecordLatency("api_"+operation, latency)
	
	if success {
		pm.RecordThroughput("api_"+operation, 1)
	} else {
		pm.RecordError("api_"+operation, 1)
	}
}

// RecordMiddlewareOperation records middleware operation metrics
func (pm *PerformanceMonitor) RecordMiddlewareOperation(operation string, latency time.Duration, success bool) {
	pm.RecordLatency("middleware_"+operation, latency)
	
	if success {
		pm.RecordThroughput("middleware_"+operation, 1)
	} else {
		pm.RecordError("middleware_"+operation, 1)
	}
}
