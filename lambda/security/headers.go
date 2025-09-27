package security

import (
	"net/http"
	"strings"
)

// SecurityHeaders provides security header management
type SecurityHeaders struct {
	// CORS settings
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	exposedHeaders []string
	allowCredentials bool
	maxAge          int
	
	// Security headers
	contentTypeOptions bool
	frameOptions      string
	xssProtection      bool
	referrerPolicy    string
	permissionsPolicy  string
	strictTransportSecurity bool
	hstsMaxAge        int
	hstsIncludeSubdomains bool
	hstsPreload       bool
}

// NewSecurityHeaders creates a new security headers manager
func NewSecurityHeaders() *SecurityHeaders {
	return &SecurityHeaders{
		allowedOrigins: []string{"*"},
		allowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		allowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		exposedHeaders: []string{},
		allowCredentials: false,
		maxAge:          86400, // 24 hours
		contentTypeOptions: true,
		frameOptions:      "DENY",
		xssProtection:      true,
		referrerPolicy:    "strict-origin-when-cross-origin",
		permissionsPolicy:  "geolocation=(), microphone=(), camera=()",
		strictTransportSecurity: true,
		hstsMaxAge:        31536000, // 1 year
		hstsIncludeSubdomains: true,
		hstsPreload:       true,
	}
}

// SetCORSOrigin sets the allowed CORS origins
func (sh *SecurityHeaders) SetCORSOrigin(origins []string) {
	sh.allowedOrigins = origins
}

// SetCORSMethods sets the allowed CORS methods
func (sh *SecurityHeaders) SetCORSMethods(methods []string) {
	sh.allowedMethods = methods
}

// SetCORSHeaders sets the allowed CORS headers
func (sh *SecurityHeaders) SetCORSHeaders(headers []string) {
	sh.allowedHeaders = headers
}

// SetCORSExposedHeaders sets the exposed CORS headers
func (sh *SecurityHeaders) SetCORSExposedHeaders(headers []string) {
	sh.exposedHeaders = headers
}

// SetCORSAllowCredentials sets whether credentials are allowed
func (sh *SecurityHeaders) SetCORSAllowCredentials(allow bool) {
	sh.allowCredentials = allow
}

// SetCORSMaxAge sets the CORS max age
func (sh *SecurityHeaders) SetCORSMaxAge(maxAge int) {
	sh.maxAge = maxAge
}

// SetContentTypeOptions sets the X-Content-Type-Options header
func (sh *SecurityHeaders) SetContentTypeOptions(enabled bool) {
	sh.contentTypeOptions = enabled
}

// SetFrameOptions sets the X-Frame-Options header
func (sh *SecurityHeaders) SetFrameOptions(options string) {
	sh.frameOptions = options
}

// SetXSSProtection sets the X-XSS-Protection header
func (sh *SecurityHeaders) SetXSSProtection(enabled bool) {
	sh.xssProtection = enabled
}

// SetReferrerPolicy sets the Referrer-Policy header
func (sh *SecurityHeaders) SetReferrerPolicy(policy string) {
	sh.referrerPolicy = policy
}

// SetPermissionsPolicy sets the Permissions-Policy header
func (sh *SecurityHeaders) SetPermissionsPolicy(policy string) {
	sh.permissionsPolicy = policy
}

// SetStrictTransportSecurity sets the Strict-Transport-Security header
func (sh *SecurityHeaders) SetStrictTransportSecurity(enabled bool, maxAge int, includeSubdomains bool, preload bool) {
	sh.strictTransportSecurity = enabled
	sh.hstsMaxAge = maxAge
	sh.hstsIncludeSubdomains = includeSubdomains
	sh.hstsPreload = preload
}

// ApplyHeaders applies security headers to an HTTP response
func (sh *SecurityHeaders) ApplyHeaders(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	sh.applyCORSHeaders(w, r)
	
	// Security headers
	sh.applySecurityHeaders(w)
	
	// Content-Type header
	w.Header().Set("Content-Type", "application/json")
}

// applyCORSHeaders applies CORS headers
func (sh *SecurityHeaders) applyCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	
	// Check if origin is allowed
	if sh.isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(sh.allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(sh.allowedHeaders, ", "))
	
	if len(sh.exposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(sh.exposedHeaders, ", "))
	}
	
	if sh.allowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	
	w.Header().Set("Access-Control-Max-Age", string(rune(sh.maxAge)))
	
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// applySecurityHeaders applies security headers
func (sh *SecurityHeaders) applySecurityHeaders(w http.ResponseWriter) {
	// X-Content-Type-Options
	if sh.contentTypeOptions {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}
	
	// X-Frame-Options
	if sh.frameOptions != "" {
		w.Header().Set("X-Frame-Options", sh.frameOptions)
	}
	
	// X-XSS-Protection
	if sh.xssProtection {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
	}
	
	// Referrer-Policy
	if sh.referrerPolicy != "" {
		w.Header().Set("Referrer-Policy", sh.referrerPolicy)
	}
	
	// Permissions-Policy
	if sh.permissionsPolicy != "" {
		w.Header().Set("Permissions-Policy", sh.permissionsPolicy)
	}
	
	// Strict-Transport-Security
	if sh.strictTransportSecurity {
		hstsValue := "max-age=" + string(rune(sh.hstsMaxAge))
		if sh.hstsIncludeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if sh.hstsPreload {
			hstsValue += "; preload"
		}
		w.Header().Set("Strict-Transport-Security", hstsValue)
	}
}

// isOriginAllowed checks if an origin is allowed
func (sh *SecurityHeaders) isOriginAllowed(origin string) bool {
	if len(sh.allowedOrigins) == 0 {
		return false
	}
	
	for _, allowedOrigin := range sh.allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	
	return false
}

// GetCORSHeaders returns CORS headers as a map
func (sh *SecurityHeaders) GetCORSHeaders() map[string]string {
	headers := make(map[string]string)
	
	headers["Access-Control-Allow-Methods"] = strings.Join(sh.allowedMethods, ", ")
	headers["Access-Control-Allow-Headers"] = strings.Join(sh.allowedHeaders, ", ")
	
	if len(sh.exposedHeaders) > 0 {
		headers["Access-Control-Expose-Headers"] = strings.Join(sh.exposedHeaders, ", ")
	}
	
	if sh.allowCredentials {
		headers["Access-Control-Allow-Credentials"] = "true"
	}
	
	headers["Access-Control-Max-Age"] = string(rune(sh.maxAge))
	
	return headers
}

// GetSecurityHeaders returns security headers as a map
func (sh *SecurityHeaders) GetSecurityHeaders() map[string]string {
	headers := make(map[string]string)
	
	if sh.contentTypeOptions {
		headers["X-Content-Type-Options"] = "nosniff"
	}
	
	if sh.frameOptions != "" {
		headers["X-Frame-Options"] = sh.frameOptions
	}
	
	if sh.xssProtection {
		headers["X-XSS-Protection"] = "1; mode=block"
	}
	
	if sh.referrerPolicy != "" {
		headers["Referrer-Policy"] = sh.referrerPolicy
	}
	
	if sh.permissionsPolicy != "" {
		headers["Permissions-Policy"] = sh.permissionsPolicy
	}
	
	if sh.strictTransportSecurity {
		hstsValue := "max-age=" + string(rune(sh.hstsMaxAge))
		if sh.hstsIncludeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if sh.hstsPreload {
			hstsValue += "; preload"
		}
		headers["Strict-Transport-Security"] = hstsValue
	}
	
	return headers
}
