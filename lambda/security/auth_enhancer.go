package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthEnhancer provides enhanced authentication features
type AuthEnhancer struct {
	// JWT settings
	secretKey     []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	
	// Rate limiting
	rateLimiter   map[string]time.Time
	maxAttempts   int
	lockoutTime   time.Duration
	
	// Password policies
	minLength     int
	requireUpper  bool
	requireLower  bool
	requireNumber bool
	requireSpecial bool
}

// NewAuthEnhancer creates a new authentication enhancer
func NewAuthEnhancer(secretKey []byte) *AuthEnhancer {
	return &AuthEnhancer{
		secretKey:     secretKey,
		tokenExpiry:   15 * time.Minute,  // Short-lived access tokens
		refreshExpiry: 7 * 24 * time.Hour, // Long-lived refresh tokens
		rateLimiter:   make(map[string]time.Time),
		maxAttempts:   5,
		lockoutTime:   15 * time.Minute,
		minLength:     8,
		requireUpper:  true,
		requireLower:  true,
		requireNumber: true,
		requireSpecial: true,
	}
}

// EnhancedUser represents a user with enhanced security features
type EnhancedUser struct {
	Username        string    `json:"username"`
	PasswordHash    string    `json:"password_hash"`
	Email           string    `json:"email"`
	IsActive        bool      `json:"is_active"`
	IsVerified      bool      `json:"is_verified"`
	LastLogin       time.Time `json:"last_login"`
	FailedAttempts  int       `json:"failed_attempts"`
	LockedUntil     *time.Time `json:"locked_until"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// ValidatePasswordStrength validates password strength
func (ae *AuthEnhancer) ValidatePasswordStrength(password string) error {
	if len(password) < ae.minLength {
		return fmt.Errorf("password must be at least %d characters long", ae.minLength)
	}
	
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= 33 && char <= 126 && !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')):
			hasSpecial = true
		}
	}
	
	if ae.requireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if ae.requireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if ae.requireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if ae.requireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

// HashPasswordWithSalt creates a salted password hash
func (ae *AuthEnhancer) HashPasswordWithSalt(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	
	// Combine password and salt
	saltedPassword := append([]byte(password), salt...)
	
	// Hash with bcrypt
	hash, err := bcrypt.GenerateFromPassword(saltedPassword, 12) // Higher cost factor
	if err != nil {
		return "", err
	}
	
	// Combine hash and salt for storage
	saltHex := hex.EncodeToString(salt)
	hashHex := hex.EncodeToString(hash)
	
	return fmt.Sprintf("%s:%s", saltHex, hashHex), nil
}

// VerifyPasswordWithSalt verifies a password against a salted hash
func (ae *AuthEnhancer) VerifyPasswordWithSalt(password, hashedPassword string) bool {
	// Split salt and hash
	parts := split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}
	
	saltHex, hashHex := parts[0], parts[1]
	
	// Decode salt and hash
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}
	
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}
	
	// Combine password and salt
	saltedPassword := append([]byte(password), salt...)
	
	// Verify password
	return bcrypt.CompareHashAndPassword(hash, saltedPassword) == nil
}

// GenerateTokenPair generates access and refresh tokens
func (ae *AuthEnhancer) GenerateTokenPair(user EnhancedUser) (*TokenPair, error) {
	// Generate access token
	accessToken, err := ae.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	
	// Generate refresh token
	refreshToken, err := ae.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ae.tokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken generates a short-lived access token
func (ae *AuthEnhancer) generateAccessToken(user EnhancedUser) (string, error) {
	claims := jwt.MapClaims{
		"user":    user.Username,
		"email":   user.Email,
		"active":  user.IsActive,
		"verified": user.IsVerified,
		"exp":     time.Now().Add(ae.tokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ae.secretKey)
}

// generateRefreshToken generates a long-lived refresh token
func (ae *AuthEnhancer) generateRefreshToken(user EnhancedUser) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	
	// Hash the token for storage
	hash := sha256.Sum256(tokenBytes)
	tokenHash := hex.EncodeToString(hash[:])
	
	// Create JWT with token hash
	claims := jwt.MapClaims{
		"user":  user.Username,
		"token": tokenHash,
		"exp":   time.Now().Add(ae.refreshExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"type":  "refresh",
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ae.secretKey)
}

// ValidateToken validates and parses a JWT token
func (ae *AuthEnhancer) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ae.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}
	
	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token has expired")
		}
	}
	
	return claims, nil
}

// CheckRateLimit checks if a user has exceeded rate limits
func (ae *AuthEnhancer) CheckRateLimit(username string) bool {
	now := time.Now()
	
	// Check if user is locked out
	if lockoutTime, exists := ae.rateLimiter[username]; exists {
		if now.Before(lockoutTime) {
			return false // Still locked out
		}
		delete(ae.rateLimiter, username) // Remove expired lockout
	}
	
	return true // Not rate limited
}

// RecordFailedAttempt records a failed login attempt
func (ae *AuthEnhancer) RecordFailedAttempt(username string) {
	ae.rateLimiter[username] = time.Now().Add(ae.lockoutTime)
}

// GenerateSecureRandom generates a secure random string
func (ae *AuthEnhancer) GenerateSecureRandom(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Helper function to split string
func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}
