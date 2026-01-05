package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Token struct {
	accessSecret string

	accessExpire  time.Duration
	refreshExpire time.Duration
}

func NewJWT(
	accessSecret string,
	accessExpire time.Duration,
	refreshExpire time.Duration,
) *Token {
	return &Token{
		accessSecret:  accessSecret,
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

func (j *Token) AccessExpire() time.Duration {
	return j.accessExpire
}

func (j *Token) RefreshExpire() time.Duration {
	return j.refreshExpire
}

func (j *Token) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()

	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.accessSecret))
}

func (j *Token) ParseAccessToken(tokenStr string) (string, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(j.accessSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil || !token.Valid || claims.UserID == "" {
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}

func (j *Token) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 64)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
