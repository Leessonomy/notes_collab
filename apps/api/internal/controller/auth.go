package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/repository"
	"notes-collab-api/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

const minPasswordLength = 8

type AuthController struct {
	userRepo         *repository.UserRepo
	refreshTokenRepo *repository.RefreshTokenRepo
	jwt              *utils.Token
	secureCookies    bool
}

func NewAuthController(
	user *repository.UserRepo,
	refreshToken *repository.RefreshTokenRepo,
	jwt *utils.Token,
	secureCookies bool,
) *AuthController {
	return &AuthController{
		userRepo:         user,
		refreshTokenRepo: refreshToken,
		jwt:              jwt,
		secureCookies:    secureCookies,
	}
}

type SignUpBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req SignUpBody) validate() error {
	if len(req.Name) == 0 {
		return errors.New("name is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email")
	}
	if len(req.Password) < minPasswordLength {
		return errors.New("password is too short")
	}
	return nil
}

func (c *AuthController) createSession(w http.ResponseWriter, r *http.Request, userID string) error {
	access, err := c.jwt.GenerateAccessToken(userID)
	if err != nil {
		return err
	}

	refresh, err := c.jwt.GenerateRefreshToken()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	accessExpiresAt := now.Add(c.jwt.AccessExpire())
	refreshExpiresAt := now.Add(c.jwt.RefreshExpire())

	err = c.refreshTokenRepo.Save(r.Context(), repository.RefreshToken{
		Token:     refresh,
		UserID:    userID,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}

	c.setSecureCookie(w, utils.AccessCookieName, access, accessExpiresAt)
	c.setSecureCookie(w, utils.RefreshCookieName, refresh, refreshExpiresAt)

	return nil
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var b SignUpBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b.Email = strings.ToLower(strings.TrimSpace(b.Email))

	if err := b.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := c.userRepo.GetByEmail(r.Context(), b.Email)
	if err == nil {
		http.Error(w, domain.ErrEmailTaken.Error(), http.StatusConflict)
		return
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		internalError(w, err)
		return
	}

	hash, err := utils.HashPassword(b.Password)
	if err != nil {
		internalError(w, err)
		return
	}

	user := domain.User{
		ID:        uuid.NewString(),
		Name:      b.Name,
		Email:     b.Email,
		Password:  hash,
		CreatedAt: time.Now().UTC(),
	}

	if err := c.userRepo.Create(r.Context(), user); err != nil {
		internalError(w, err)
		return
	}

	if err := c.createSession(w, r, user.ID); err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (c *AuthController) LogIn(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	b.Email = strings.ToLower(strings.TrimSpace(b.Email))

	user, err := c.userRepo.GetByEmail(r.Context(), b.Email)

	if err != nil {
		http.Error(w, domain.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(user.Password, b.Password) {
		http.Error(w, domain.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		return
	}

	if err := c.createSession(w, r, user.ID); err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.RefreshCookieName)
	if err != nil {
		unauthorized(w)
		return
	}

	newRefresh, err := c.jwt.GenerateRefreshToken()
	if err != nil {
		internalError(w, err)
		return
	}

	now := time.Now().UTC()
	refreshExpiresAt := now.Add(c.jwt.RefreshExpire())

	userID, err := c.refreshTokenRepo.Rotate(r.Context(), cookie.Value, newRefresh, refreshExpiresAt)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			c.clearCookies(w, utils.AccessCookieName)
			c.clearCookies(w, utils.RefreshCookieName)
			unauthorized(w)
			return
		}
		internalError(w, err)
		return
	}

	access, err := c.jwt.GenerateAccessToken(userID)
	if err != nil {
		internalError(w, err)
		return
	}

	c.setSecureCookie(w, utils.AccessCookieName, access, now.Add(c.jwt.AccessExpire()))
	c.setSecureCookie(w, utils.RefreshCookieName, newRefresh, refreshExpiresAt)

	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) LogOut(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(utils.RefreshCookieName); err == nil {
		c.refreshTokenRepo.Delete(r.Context(), cookie.Value)
	}

	c.clearCookies(w, utils.AccessCookieName)
	c.clearCookies(w, utils.RefreshCookieName)
	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) setSecureCookie(w http.ResponseWriter, name, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AuthController) clearCookies(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
