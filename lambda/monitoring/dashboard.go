package monitoring

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// DashboardConfig provides CloudWatch dashboard configuration
type DashboardConfig struct {
	cloudwatchClient *cloudwatch.Client
	dashboardName    string
	namespace        string
}

// NewDashboardConfig creates a new dashboard configuration
func NewDashboardConfig(cloudwatchClient *cloudwatch.Client, dashboardName, namespace string) *DashboardConfig {
	return &DashboardConfig{
		cloudwatchClient: cloudwatchClient,
		dashboardName:    dashboardName,
		namespace:        namespace,
	}
}

// CreateDashboard creates a CloudWatch dashboard
func (dc *DashboardConfig) CreateDashboard(ctx context.Context) error {
	// Create dashboard body
	dashboardBody := dc.generateDashboardBody()

	// Create dashboard
	_, err := dc.cloudwatchClient.PutDashboard(ctx, &cloudwatch.PutDashboardInput{
		DashboardName: &dc.dashboardName,
		DashboardBody: &dashboardBody,
	})

	return err
}

// generateDashboardBody generates the dashboard body
func (dc *DashboardConfig) generateDashboardBody() string {
	dashboard := map[string]interface{}{
		"widgets": []map[string]interface{}{
			// Request metrics widget
			{
				"type":   "metric",
				"x":      0,
				"y":      0,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "requests_total"},
						{dc.namespace, "requests_by_status"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Request Metrics",
					"period":  300,
				},
			},
			// Response time widget
			{
				"type":   "metric",
				"x":      12,
				"y":      0,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "request_duration_seconds"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Response Time",
					"period":  300,
					"stat":    "Average",
				},
			},
			// Database metrics widget
			{
				"type":   "metric",
				"x":      0,
				"y":      6,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "database_operations_total"},
						{dc.namespace, "database_operations_success"},
						{dc.namespace, "database_operations_failure"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Database Operations",
					"period":  300,
				},
			},
			// Authentication metrics widget
			{
				"type":   "metric",
				"x":      12,
				"y":      6,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "authentication_attempts_total"},
						{dc.namespace, "authentication_success"},
						{dc.namespace, "authentication_failure"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Authentication Metrics",
					"period":  300,
				},
			},
			// Error metrics widget
			{
				"type":   "metric",
				"x":      0,
				"y":      12,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "errors_total"},
						{dc.namespace, "errors_by_type"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Error Metrics",
					"period":  300,
				},
			},
			// Security metrics widget
			{
				"type":   "metric",
				"x":      12,
				"y":      12,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "security_events_total"},
						{dc.namespace, "security_events_by_type"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Security Metrics",
					"period":  300,
				},
			},
			// Memory usage widget
			{
				"type":   "metric",
				"x":      0,
				"y":      18,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "memory_usage_bytes"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Memory Usage",
					"period":  300,
					"stat":    "Average",
				},
			},
			// Performance metrics widget
			{
				"type":   "metric",
				"x":      12,
				"y":      18,
				"width":  12,
				"height": 6,
				"properties": map[string]interface{}{
					"metrics": [][]interface{}{
						{dc.namespace, "operation_duration_seconds"},
					},
					"view":    "timeSeries",
					"stacked": false,
					"region":  "us-east-1",
					"title":   "Operation Duration",
					"period":  300,
					"stat":    "Average",
				},
			},
		},
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(dashboard)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// CreateAlarms creates CloudWatch alarms
func (dc *DashboardConfig) CreateAlarms(ctx context.Context) error {
	alarms := []types.PutMetricAlarmInput{
		// High error rate alarm
		{
			AlarmName:          &[]string{"HighErrorRate"}[0],
			AlarmDescription:   &[]string{"High error rate detected"}[0],
			MetricName:         &[]string{"errors_total"}[0],
			Namespace:          &dc.namespace,
			Statistic:          types.StatisticSum,
			Period:             &[]int32{300}[0],
			EvaluationPeriods:  &[]int32{2}[0],
			Threshold:          &[]float64{10}[0],
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			AlarmActions:       []string{},
		},
		// High response time alarm
		{
			AlarmName:          &[]string{"HighResponseTime"}[0],
			AlarmDescription:   &[]string{"High response time detected"}[0],
			MetricName:         &[]string{"request_duration_seconds"}[0],
			Namespace:          &dc.namespace,
			Statistic:          types.StatisticAverage,
			Period:             &[]int32{300}[0],
			EvaluationPeriods:  &[]int32{2}[0],
			Threshold:          &[]float64{5}[0],
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			AlarmActions:       []string{},
		},
		// High memory usage alarm
		{
			AlarmName:          &[]string{"HighMemoryUsage"}[0],
			AlarmDescription:   &[]string{"High memory usage detected"}[0],
			MetricName:         &[]string{"memory_usage_bytes"}[0],
			Namespace:          &dc.namespace,
			Statistic:          types.StatisticAverage,
			Period:             &[]int32{300}[0],
			EvaluationPeriods:  &[]int32{2}[0],
			Threshold:          &[]float64{100000000}[0], // 100MB
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			AlarmActions:       []string{},
		},
		// Authentication failure alarm
		{
			AlarmName:          &[]string{"HighAuthenticationFailures"}[0],
			AlarmDescription:   &[]string{"High authentication failure rate detected"}[0],
			MetricName:         &[]string{"authentication_failure"}[0],
			Namespace:          &dc.namespace,
			Statistic:          types.StatisticSum,
			Period:             &[]int32{300}[0],
			EvaluationPeriods:  &[]int32{2}[0],
			Threshold:          &[]float64{5}[0],
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			AlarmActions:       []string{},
		},
		// Security events alarm
		{
			AlarmName:          &[]string{"SecurityEvents"}[0],
			AlarmDescription:   &[]string{"Security events detected"}[0],
			MetricName:         &[]string{"security_events_total"}[0],
			Namespace:          &dc.namespace,
			Statistic:          types.StatisticSum,
			Period:             &[]int32{300}[0],
			EvaluationPeriods:  &[]int32{1}[0],
			Threshold:          &[]float64{1}[0],
			ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
			AlarmActions:       []string{},
		},
	}

	// Create each alarm
	for _, alarm := range alarms {
		_, err := dc.cloudwatchClient.PutMetricAlarm(ctx, &alarm)
		if err != nil {
			return fmt.Errorf("failed to create alarm %s: %w", *alarm.AlarmName, err)
		}
	}

	return nil
}

// GetDashboardURL returns the dashboard URL
func (dc *DashboardConfig) GetDashboardURL() string {
	return fmt.Sprintf("https://console.aws.amazon.com/cloudwatch/home?region=us-east-1#dashboards:name=%s", dc.dashboardName)
}

// UpdateDashboard updates the dashboard
func (dc *DashboardConfig) UpdateDashboard(ctx context.Context) error {
	return dc.CreateDashboard(ctx)
}

// DeleteDashboard deletes the dashboard
func (dc *DashboardConfig) DeleteDashboard(ctx context.Context) error {
	_, err := dc.cloudwatchClient.DeleteDashboards(ctx, &cloudwatch.DeleteDashboardsInput{
		DashboardNames: []string{dc.dashboardName},
	})

	return err
}
