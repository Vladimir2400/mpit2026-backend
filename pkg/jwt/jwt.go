package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token creation and validation
type TokenManager struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(accessSecret, refreshSecret string, accessExpiryMin, refreshExpiryDay int) *TokenManager {
	return &TokenManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  time.Duration(accessExpiryMin) * time.Minute,
		refreshExpiry: time.Duration(refreshExpiryDay) * 24 * time.Hour,
	}
}

// GenerateAccessToken generates a new access token
func (tm *TokenManager) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.accessSecret))
}

// GenerateRefreshToken generates a new refresh token
func (tm *TokenManager) GenerateRefreshToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.refreshSecret))
}

// ValidateAccessToken validates an access token and returns claims
func (tm *TokenManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return tm.validateToken(tokenString, tm.accessSecret)
}

// ValidateRefreshToken validates a refresh token and returns claims
func (tm *TokenManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return tm.validateToken(tokenString, tm.refreshSecret)
}

// validateToken validates a token with the given secret
func (tm *TokenManager) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractUserID extracts user ID from access token without validation
// Use only for non-critical operations
func (tm *TokenManager) ExtractUserID(tokenString string) (uuid.UUID, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return claims.UserID, nil
}
