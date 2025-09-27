package security

import (
	"context"
	"sync"
	"time"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	// Rate limiting algorithms
	tokenBucket map[string]*TokenBucket
	slidingWindow map[string]*SlidingWindow
	
	// Configuration
	maxRequests    int
	windowSize     time.Duration
	cleanupInterval time.Duration
	
	// Mutex for thread safety
	mu sync.RWMutex
	
	// Cleanup ticker
	cleanupTicker *time.Ticker
	stopChan     chan bool
}

// TokenBucket implements token bucket algorithm
type TokenBucket struct {
	Capacity     int
	Tokens       int
	LastRefill   time.Time
	RefillRate   int
	RefillPeriod time.Duration
}

// SlidingWindow implements sliding window algorithm
type SlidingWindow struct {
	Requests []time.Time
	Window   time.Duration
	MaxRequests int
}

// RateLimitResult represents the result of rate limiting check
type RateLimitResult struct {
	Allowed    bool
	Remaining  int
	ResetTime  time.Time
	RetryAfter time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, windowSize time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokenBucket:    make(map[string]*TokenBucket),
		slidingWindow:  make(map[string]*SlidingWindow),
		maxRequests:    maxRequests,
		windowSize:     windowSize,
		cleanupInterval: 5 * time.Minute,
		stopChan:       make(chan bool),
	}
	
	// Start cleanup routine
	rl.startCleanup()
	
	return rl
}

// CheckRateLimit checks if a request is allowed based on rate limiting
func (rl *RateLimiter) CheckRateLimit(identifier string, algorithm string) RateLimitResult {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	switch algorithm {
	case "token_bucket":
		return rl.checkTokenBucket(identifier)
	case "sliding_window":
		return rl.checkSlidingWindow(identifier)
	default:
		return rl.checkSlidingWindow(identifier) // Default to sliding window
	}
}

// checkTokenBucket checks rate limit using token bucket algorithm
func (rl *RateLimiter) checkTokenBucket(identifier string) RateLimitResult {
	bucket, exists := rl.tokenBucket[identifier]
	if !exists {
		bucket = &TokenBucket{
			Capacity:     rl.maxRequests,
			Tokens:       rl.maxRequests,
			LastRefill:   time.Now(),
			RefillRate:   1,
			RefillPeriod: time.Second,
		}
		rl.tokenBucket[identifier] = bucket
	}
	
	now := time.Now()
	
	// Refill tokens
	timeSinceLastRefill := now.Sub(bucket.LastRefill)
	tokensToAdd := int(timeSinceLastRefill / bucket.RefillPeriod) * bucket.RefillRate
	
	if tokensToAdd > 0 {
		bucket.Tokens = min(bucket.Capacity, bucket.Tokens+tokensToAdd)
		bucket.LastRefill = now
	}
	
	// Check if request is allowed
	if bucket.Tokens > 0 {
		bucket.Tokens--
		return RateLimitResult{
			Allowed:   true,
			Remaining: bucket.Tokens,
			ResetTime: now.Add(bucket.RefillPeriod),
		}
	}
	
	// Calculate retry after
	retryAfter := bucket.RefillPeriod - timeSinceLastRefill
	if retryAfter < 0 {
		retryAfter = 0
	}
	
	return RateLimitResult{
		Allowed:    false,
		Remaining: 0,
		ResetTime:  now.Add(retryAfter),
		RetryAfter: retryAfter,
	}
}

// checkSlidingWindow checks rate limit using sliding window algorithm
func (rl *RateLimiter) checkSlidingWindow(identifier string) RateLimitResult {
	window, exists := rl.slidingWindow[identifier]
	if !exists {
		window = &SlidingWindow{
			Requests:    make([]time.Time, 0),
			Window:      rl.windowSize,
			MaxRequests: rl.maxRequests,
		}
		rl.slidingWindow[identifier] = window
	}
	
	now := time.Now()
	
	// Remove old requests outside the window
	cutoff := now.Add(-window.Window)
	var validRequests []time.Time
	for _, requestTime := range window.Requests {
		if requestTime.After(cutoff) {
			validRequests = append(validRequests, requestTime)
		}
	}
	window.Requests = validRequests
	
	// Check if request is allowed
	if len(window.Requests) < window.MaxRequests {
		window.Requests = append(window.Requests, now)
		return RateLimitResult{
			Allowed:   true,
			Remaining: window.MaxRequests - len(window.Requests),
			ResetTime: now.Add(window.Window),
		}
	}
	
	// Calculate retry after (time until oldest request expires)
	oldestRequest := window.Requests[0]
	retryAfter := oldestRequest.Add(window.Window).Sub(now)
	
	return RateLimitResult{
		Allowed:    false,
		Remaining:  0,
		ResetTime:  now.Add(retryAfter),
		RetryAfter: retryAfter,
	}
}

// startCleanup starts the cleanup routine
func (rl *RateLimiter) startCleanup() {
	rl.cleanupTicker = time.NewTicker(rl.cleanupInterval)
	
	go func() {
		for {
			select {
			case <-rl.cleanupTicker.C:
				rl.cleanup()
			case <-rl.stopChan:
				rl.cleanupTicker.Stop()
				return
			}
		}
	}()
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.windowSize * 2) // Keep entries for 2x window size
	
	// Cleanup token buckets
	for identifier, bucket := range rl.tokenBucket {
		if bucket.LastRefill.Before(cutoff) {
			delete(rl.tokenBucket, identifier)
		}
	}
	
	// Cleanup sliding windows
	for identifier, window := range rl.slidingWindow {
		if len(window.Requests) == 0 || window.Requests[len(window.Requests)-1].Before(cutoff) {
			delete(rl.slidingWindow, identifier)
		}
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	rl.stopChan <- true
}

// GetStats returns rate limiter statistics
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	return map[string]interface{}{
		"token_buckets":   len(rl.tokenBucket),
		"sliding_windows": len(rl.slidingWindow),
		"max_requests":    rl.maxRequests,
		"window_size":     rl.windowSize.String(),
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IPRateLimiter provides IP-based rate limiting
type IPRateLimiter struct {
	rateLimiter *RateLimiter
	blockedIPs  map[string]time.Time
	mu          sync.RWMutex
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(maxRequests int, windowSize time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		rateLimiter: NewRateLimiter(maxRequests, windowSize),
		blockedIPs:  make(map[string]time.Time),
	}
}

// CheckIPRateLimit checks rate limit for an IP address
func (irl *IPRateLimiter) CheckIPRateLimit(ip string) RateLimitResult {
	// Check if IP is blocked
	irl.mu.RLock()
	if blockTime, exists := irl.blockedIPs[ip]; exists {
		if time.Now().Before(blockTime) {
			irl.mu.RUnlock()
			return RateLimitResult{
				Allowed:    false,
				Remaining:  0,
				RetryAfter: blockTime.Sub(time.Now()),
			}
		}
		// Remove expired block
		delete(irl.blockedIPs, ip)
	}
	irl.mu.RUnlock()
	
	// Check rate limit
	result := irl.rateLimiter.CheckRateLimit(ip, "sliding_window")
	
	// Block IP if too many requests
	if !result.Allowed && result.RetryAfter > 0 {
		irl.mu.Lock()
		irl.blockedIPs[ip] = time.Now().Add(result.RetryAfter)
		irl.mu.Unlock()
	}
	
	return result
}

// UnblockIP unblocks an IP address
func (irl *IPRateLimiter) UnblockIP(ip string) {
	irl.mu.Lock()
	defer irl.mu.Unlock()
	
	delete(irl.blockedIPs, ip)
}

// GetBlockedIPs returns list of blocked IPs
func (irl *IPRateLimiter) GetBlockedIPs() []string {
	irl.mu.RLock()
	defer irl.mu.RUnlock()
	
	var blockedIPs []string
	for ip, blockTime := range irl.blockedIPs {
		if time.Now().Before(blockTime) {
			blockedIPs = append(blockedIPs, ip)
		}
	}
	
	return blockedIPs
}
