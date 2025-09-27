package security

import (
	"regexp"
	"strings"
	"unicode"
)

// InputValidator provides comprehensive input validation
type InputValidator struct {
	// Regex patterns for validation
	usernamePattern *regexp.Regexp
	emailPattern    *regexp.Regexp
	phonePattern    *regexp.Regexp
	
	// XSS patterns
	xssPatterns []*regexp.Regexp
	
	// SQL injection patterns
	sqlInjectionPatterns []*regexp.Regexp
}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	return &InputValidator{
		usernamePattern: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`),
		emailPattern:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		phonePattern:    regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
		xssPatterns: []*regexp.Regexp{
			regexp.MustCompile(`<script[^>]*>.*?</script>`),
			regexp.MustCompile(`javascript:`),
			regexp.MustCompile(`on\w+\s*=`),
			regexp.MustCompile(`<iframe[^>]*>.*?</iframe>`),
			regexp.MustCompile(`<object[^>]*>.*?</object>`),
			regexp.MustCompile(`<embed[^>]*>.*?</embed>`),
		},
		sqlInjectionPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`),
			regexp.MustCompile(`(?i)(or|and)\s+\d+\s*=\s*\d+`),
			regexp.MustCompile(`(?i)(or|and)\s+'.*'\s*=\s*'.*'`),
			regexp.MustCompile(`(?i)(or|and)\s+".*"\s*=\s*".*"`),
			regexp.MustCompile(`(?i)(or|and)\s+\w+\s*=\s*\w+`),
			regexp.MustCompile(`(?i)(or|and)\s+\w+\s*like\s*'%`),
			regexp.MustCompile(`(?i)(or|and)\s+\w+\s*like\s*"%`),
		},
	}
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// ValidateUsername validates username input
func (iv *InputValidator) ValidateUsername(username string) ValidationResult {
	var errors []string
	
	// Check length
	if len(username) < 3 {
		errors = append(errors, "username must be at least 3 characters long")
	}
	if len(username) > 20 {
		errors = append(errors, "username must be no more than 20 characters long")
	}
	
	// Check pattern
	if !iv.usernamePattern.MatchString(username) {
		errors = append(errors, "username can only contain letters, numbers, underscores, and hyphens")
	}
	
	// Check for XSS
	if iv.containsXSS(username) {
		errors = append(errors, "username contains potentially malicious content")
	}
	
	// Check for SQL injection
	if iv.containsSQLInjection(username) {
		errors = append(errors, "username contains potentially malicious SQL content")
	}
	
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateEmail validates email input
func (iv *InputValidator) ValidateEmail(email string) ValidationResult {
	var errors []string
	
	// Check length
	if len(email) > 254 {
		errors = append(errors, "email must be no more than 254 characters long")
	}
	
	// Check pattern
	if !iv.emailPattern.MatchString(email) {
		errors = append(errors, "email format is invalid")
	}
	
	// Check for XSS
	if iv.containsXSS(email) {
		errors = append(errors, "email contains potentially malicious content")
	}
	
	// Check for SQL injection
	if iv.containsSQLInjection(email) {
		errors = append(errors, "email contains potentially malicious SQL content")
	}
	
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidatePassword validates password input
func (iv *InputValidator) ValidatePassword(password string) ValidationResult {
	var errors []string
	
	// Check length
	if len(password) < 8 {
		errors = append(errors, "password must be at least 8 characters long")
	}
	if len(password) > 128 {
		errors = append(errors, "password must be no more than 128 characters long")
	}
	
	// Check for common weak passwords
	if iv.isWeakPassword(password) {
		errors = append(errors, "password is too weak")
	}
	
	// Check for XSS
	if iv.containsXSS(password) {
		errors = append(errors, "password contains potentially malicious content")
	}
	
	// Check for SQL injection
	if iv.containsSQLInjection(password) {
		errors = append(errors, "password contains potentially malicious SQL content")
	}
	
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidatePhone validates phone number input
func (iv *InputValidator) ValidatePhone(phone string) ValidationResult {
	var errors []string
	
	// Check pattern
	if !iv.phonePattern.MatchString(phone) {
		errors = append(errors, "phone number format is invalid")
	}
	
	// Check for XSS
	if iv.containsXSS(phone) {
		errors = append(errors, "phone number contains potentially malicious content")
	}
	
	// Check for SQL injection
	if iv.containsSQLInjection(phone) {
		errors = append(errors, "phone number contains potentially malicious SQL content")
	}
	
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateGeneralInput validates general text input
func (iv *InputValidator) ValidateGeneralInput(input string, maxLength int) ValidationResult {
	var errors []string
	
	// Check length
	if len(input) > maxLength {
		errors = append(errors, "input exceeds maximum length")
	}
	
	// Check for XSS
	if iv.containsXSS(input) {
		errors = append(errors, "input contains potentially malicious content")
	}
	
	// Check for SQL injection
	if iv.containsSQLInjection(input) {
		errors = append(errors, "input contains potentially malicious SQL content")
	}
	
	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// SanitizeInput sanitizes input by removing potentially dangerous content
func (iv *InputValidator) SanitizeInput(input string) string {
	// Remove XSS patterns
	sanitized := input
	for _, pattern := range iv.xssPatterns {
		sanitized = pattern.ReplaceAllString(sanitized, "")
	}
	
	// Remove SQL injection patterns
	for _, pattern := range iv.sqlInjectionPatterns {
		sanitized = pattern.ReplaceAllString(sanitized, "")
	}
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// containsXSS checks if input contains XSS patterns
func (iv *InputValidator) containsXSS(input string) bool {
	for _, pattern := range iv.xssPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// containsSQLInjection checks if input contains SQL injection patterns
func (iv *InputValidator) containsSQLInjection(input string) bool {
	for _, pattern := range iv.sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// isWeakPassword checks if password is weak
func (iv *InputValidator) isWeakPassword(password string) bool {
	// Common weak passwords
	weakPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
		"1234567890", "password1", "qwerty123", "dragon", "master",
	}
	
	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if lowerPassword == weak {
			return true
		}
	}
	
	// Check for sequential characters
	if iv.hasSequentialChars(password) {
		return true
	}
	
	// Check for repeated characters
	if iv.hasRepeatedChars(password) {
		return true
	}
	
	return false
}

// hasSequentialChars checks for sequential characters
func (iv *InputValidator) hasSequentialChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if unicode.IsLetter(rune(password[i])) && unicode.IsLetter(rune(password[i+1])) && unicode.IsLetter(rune(password[i+2])) {
			if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
				return true
			}
		}
	}
	return false
}

// hasRepeatedChars checks for repeated characters
func (iv *InputValidator) hasRepeatedChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}
