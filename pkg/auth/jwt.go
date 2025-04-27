package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid            = errors.New("token is invalid")
	ErrTokenExpired            = errors.New("token has expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrMissingClaims           = errors.New("required claims are missing")
)

type JWTManager struct {
	accessSecretKey      string
	refreshSecretKey     string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

type AuthClaims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func NewJWTManager(accessSecret, refreshSecret, issuer string, accessDuration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		accessSecretKey:      accessSecret,
		refreshSecretKey:     refreshSecret,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
		issuer:               issuer,
	}
}

func (m *JWTManager) GenerateTokens(phoneNumber string) (accessToken string, refreshToken string, err error) {
	accessClaims := AuthClaims{
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   phoneNumber,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.accessSecretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := AuthClaims{
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   phoneNumber,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.refreshSecretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) ValidateAccessToken(tokenString string) (*AuthClaims, error) {
	return m.validateToken(tokenString, m.accessSecretKey)
}

func (m *JWTManager) ValidateRefreshToken(tokenString string) (*AuthClaims, error) {
	return m.validateToken(tokenString, m.refreshSecretKey)
}

func (m *JWTManager) validateToken(tokenString string, secretKey string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrTokenInvalid
		}
		return nil, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		if claims.Issuer != m.issuer {
			return nil, fmt.Errorf("%w: invalid issuer", ErrTokenInvalid)
		}
		if claims.PhoneNumber == "" {
			return nil, ErrMissingClaims
		}
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
