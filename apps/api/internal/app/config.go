package app

import (
	"errors"
	"os"
	"time"
)

type config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	CORSOrigin    string
	CookieSecure  bool
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

func loadConfig() (*config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return &config{
		Port:          envOr("PORT", "8080"),
		DatabaseURL:   databaseURL,
		JWTSecret:     jwtSecret,
		CORSOrigin:    envOr("CORS_ORIGIN", "http://localhost:4200"),
		CookieSecure:  os.Getenv("COOKIE_SECURE") == "true",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 30 * 24 * time.Hour,
	}, nil
}

func envOr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
