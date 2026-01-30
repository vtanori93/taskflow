package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Claims estructura para los claims del JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

// Manager maneja la creación y validación de JWTs
type Manager struct {
	secret            string
	expirationTime    int64
	refreshExpiration int64
}

// NewManager crea un nuevo JWT Manager
func NewManager(secret string, expirationTime, refreshExpiration int64) *Manager {
	return &Manager{
		secret:            secret,
		expirationTime:    expirationTime,
		refreshExpiration: refreshExpiration,
	}
}

// GenerateToken genera un nuevo access token
func (m *Manager) GenerateToken(userID, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.expirationTime) * time.Second)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expirationTime),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("error al generar token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken genera un nuevo refresh token
func (m *Manager) GenerateRefreshToken(userID, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.refreshExpiration) * time.Second)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expirationTime),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("error al generar refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un token y retorna los claims
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (interface{}, error) {
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}
