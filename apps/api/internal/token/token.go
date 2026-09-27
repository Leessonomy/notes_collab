package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"notes-collab-api/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalid = errors.New("invalid token")

type userClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret string

	accessExpire  time.Duration
	refreshExpire time.Duration
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		accessSecret:  cfg.JWTSecret,
		accessExpire:  cfg.AccessExpire,
		refreshExpire: cfg.RefreshExpire,
	}
}

func (m *Manager) AccessExpire() time.Duration {
	return m.accessExpire
}

func (m *Manager) RefreshExpire() time.Duration {
	return m.refreshExpire
}

func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()

	claims := userClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.accessSecret))
}

func (m *Manager) ParseAccessToken(tokenStr string) (string, error) {
	claims := &userClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(m.accessSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil || !token.Valid || claims.UserID == "" {
		return "", ErrInvalid
	}

	return claims.UserID, nil
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 64)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
