package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
	"strings"
	"time"
)

const minPasswordLength = 8

type AuthController struct {
	authUseCase   *usecase.Auth
	secureCookies bool
}

func NewAuthController(authUc *usecase.Auth, secureCookies bool) *AuthController {
	return &AuthController{
		authUseCase:   authUc,
		secureCookies: secureCookies,
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

func (c *AuthController) setSession(w http.ResponseWriter, session dto.SessionOutput) {
	c.setSecureCookie(w, utils.AccessCookieName, session.AccessToken, session.AccessExpiresAt)
	c.setSecureCookie(w, utils.RefreshCookieName, session.RefreshToken, session.RefreshExpiresAt)
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(utils.RefreshCookieName)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	userID := utils.UserIDFromContext(r.Context())

	user, _ := c.authUseCase.GetUserSession(r.Context(), userID)

	writeJSON(w, http.StatusOK, user)

}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var body SignUpBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid body")
		return
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	if err := body.validate(); err != nil {
		badRequest(w, err.Error())
		return
	}

	session, err := c.authUseCase.SignUp(r.Context(), dto.SignUpInput{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	c.setSession(w, session)

	writeJSON(w, http.StatusCreated, session.User)
}

func (c *AuthController) LogIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid body")
		return
	}

	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	session, err := c.authUseCase.LogIn(r.Context(), dto.LogInInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	c.setSession(w, session)

	writeJSON(w, http.StatusOK, session.User)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.RefreshCookieName)
	if err != nil {
		unauthorized(w)
		return
	}

	session, err := c.authUseCase.Refresh(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			c.clearCookies(w, utils.AccessCookieName)
			c.clearCookies(w, utils.RefreshCookieName)
		}
		domainError(w, err)
		return
	}

	c.setSession(w, session)

	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) LogOut(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(utils.RefreshCookieName); err == nil {
		c.authUseCase.LogOut(r.Context(), cookie.Value)
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
