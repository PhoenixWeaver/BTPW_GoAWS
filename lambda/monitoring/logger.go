package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of LogLevel
func (ll LogLevel) String() string {
	switch ll {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging capabilities
type Logger struct {
	// AWS CloudWatch Logs client
	cloudwatchClient *cloudwatchlogs.Client
	
	// Log group and stream
	logGroup  string
	logStream string
	
	// Log level
	level LogLevel
	
	// Context
	ctx context.Context
	
	// Structured logging
	fields map[string]interface{}
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	Service     string                 `json:"service"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	TraceID     string                 `json:"trace_id,omitempty"`
	SpanID      string                 `json:"span_id,omitempty"`
}

// NewLogger creates a new logger
func NewLogger(ctx context.Context, cloudwatchClient *cloudwatchlogs.Client, logGroup, logStream string) *Logger {
	return &Logger{
		cloudwatchClient: cloudwatchClient,
		logGroup:         logGroup,
		logStream:        logStream,
		level:            INFO,
		ctx:              ctx,
		fields:           make(map[string]interface{}),
	}
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetFields sets additional fields for structured logging
func (l *Logger) SetFields(fields map[string]interface{}) {
	l.fields = fields
}

// AddField adds a field to the logger
func (l *Logger) AddField(key string, value interface{}) {
	l.fields[key] = value
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, err error, fields ...map[string]interface{}) {
	errorFields := map[string]interface{}{
		"error": err.Error(),
	}
	
	// Add stack trace for errors
	if pc, file, line, ok := runtime.Caller(1); ok {
		errorFields["stack_trace"] = fmt.Sprintf("%s:%d %s", file, line, runtime.FuncForPC(pc).Name())
	}
	
	// Merge with additional fields
	for _, f := range fields {
		for k, v := range f {
			errorFields[k] = v
		}
	}
	
	l.log(ERROR, message, errorFields)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	l.log(FATAL, message, fields...)
	os.Exit(1)
}

// log logs a message with the specified level
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	if level < l.level {
		return
	}
	
	// Create log entry
	entry := LogEntry{
		Timestamp:   time.Now(),
		Level:       level.String(),
		Message:     message,
		Service:     "lambda-function",
		Version:     "1.0.0",
		Environment: os.Getenv("ENVIRONMENT"),
		Fields:      make(map[string]interface{}),
	}
	
	// Add logger fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	
	// Add additional fields
	for _, f := range fields {
		for k, v := range f {
			entry.Fields[k] = v
		}
	}
	
	// Log to console
	l.logToConsole(entry)
	
	// Log to CloudWatch
	if l.cloudwatchClient != nil {
		l.logToCloudWatch(entry)
	}
}

// logToConsole logs to console
func (l *Logger) logToConsole(entry LogEntry) {
	// Format log entry
	formatted := l.formatLogEntry(entry)
	
	// Log to appropriate level
	switch entry.Level {
	case "DEBUG":
		log.Printf("[DEBUG] %s", formatted)
	case "INFO":
		log.Printf("[INFO] %s", formatted)
	case "WARN":
		log.Printf("[WARN] %s", formatted)
	case "ERROR":
		log.Printf("[ERROR] %s", formatted)
	case "FATAL":
		log.Printf("[FATAL] %s", formatted)
	}
}

// logToCloudWatch logs to CloudWatch
func (l *Logger) logToCloudWatch(entry LogEntry) {
	// Format log entry
	formatted := l.formatLogEntry(entry)
	
	// Create log event
	logEvent := types.InputLogEvent{
		Message:   &formatted,
		Timestamp: &[]int64{entry.Timestamp.UnixMilli()}[0],
	}
	
	// Send to CloudWatch
	_, err := l.cloudwatchClient.PutLogEvents(l.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &l.logGroup,
		LogStreamName: &l.logStream,
		LogEvents:     []types.InputLogEvent{logEvent},
	})
	
	if err != nil {
		// Fallback to console if CloudWatch fails
		log.Printf("[ERROR] Failed to send log to CloudWatch: %v", err)
	}
}

// formatLogEntry formats a log entry
func (l *Logger) formatLogEntry(entry LogEntry) string {
	// Convert to JSON
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf("Failed to marshal log entry: %v", err)
	}
	
	return string(jsonBytes)
}

// GetLogger returns a logger instance
func GetLogger(ctx context.Context) *Logger {
	// Create CloudWatch client
	cloudwatchClient := cloudwatchlogs.NewFromConfig(nil)
	
	// Get log group and stream from environment
	logGroup := os.Getenv("LOG_GROUP")
	if logGroup == "" {
		logGroup = "/aws/lambda/function"
	}
	
	logStream := os.Getenv("LOG_STREAM")
	if logStream == "" {
		logStream = "lambda-function"
	}
	
	return NewLogger(ctx, cloudwatchClient, logGroup, logStream)
}

// RequestLogger provides request-specific logging
type RequestLogger struct {
	*Logger
	requestID string
	userID    string
	traceID   string
	spanID    string
}

// NewRequestLogger creates a new request logger
func NewRequestLogger(ctx context.Context, requestID, userID, traceID, spanID string) *RequestLogger {
	logger := GetLogger(ctx)
	
	// Set request-specific fields
	logger.SetFields(map[string]interface{}{
		"request_id": requestID,
		"user_id":    userID,
		"trace_id":   traceID,
		"span_id":    spanID,
	})
	
	return &RequestLogger{
		Logger:    logger,
		requestID: requestID,
		userID:    userID,
		traceID:   traceID,
		spanID:    spanID,
	}
}

// LogRequestStart logs the start of a request
func (rl *RequestLogger) LogRequestStart(method, path string) {
	rl.Info("Request started", map[string]interface{}{
		"method": method,
		"path":   path,
		"action": "request_start",
	})
}

// LogRequestEnd logs the end of a request
func (rl *RequestLogger) LogRequestEnd(method, path string, statusCode int, duration time.Duration) {
	rl.Info("Request completed", map[string]interface{}{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.String(),
		"action":      "request_end",
	})
}

// LogDatabaseOperation logs a database operation
func (rl *RequestLogger) LogDatabaseOperation(operation string, table string, duration time.Duration, success bool) {
	level := INFO
	if !success {
		level = ERROR
	}
	
	rl.log(level, "Database operation", map[string]interface{}{
		"operation": operation,
		"table":     table,
		"duration":  duration.String(),
		"success":   success,
		"action":    "database_operation",
	})
}

// LogAuthentication logs authentication events
func (rl *RequestLogger) LogAuthentication(username string, success bool, reason string) {
	level := INFO
	if !success {
		level = WARN
	}
	
	rl.log(level, "Authentication attempt", map[string]interface{}{
		"username": username,
		"success":  success,
		"reason":   reason,
		"action":   "authentication",
	})
}

// LogSecurityEvent logs security events
func (rl *RequestLogger) LogSecurityEvent(event string, details map[string]interface{}) {
	rl.Warn("Security event", map[string]interface{}{
		"event":   event,
		"details": details,
		"action":  "security_event",
	})
}
