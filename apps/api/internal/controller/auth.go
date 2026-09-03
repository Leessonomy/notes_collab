package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
	"time"
)

type AuthController struct {
	authUseCase   *usecase.Auth
	secureCookies bool
}

func NewAuthController(authUC *usecase.Auth, isSecure bool) *AuthController {
	return &AuthController{
		authUseCase:   authUC,
		secureCookies: isSecure,
	}
}

func (c *AuthController) setSession(w http.ResponseWriter, session dto.SessionOutput) {
	c.setSecureCookie(w, utils.AccessCookieName, session.AccessToken, session.AccessExpiresAt)
	c.setSecureCookie(w, utils.RefreshCookieName, session.RefreshToken, session.RefreshExpiresAt)
}

func (c *AuthController) GetMe(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(utils.RefreshCookieName)
	if err != nil {
		utils.Unauthorized(w)
		return
	}

	userID := utils.UserIDFromContext(r.Context())

	user, err := c.authUseCase.GetUserSession(r.Context(), userID)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.SignUpInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	session, err := c.authUseCase.SignUp(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	c.setSession(w, session)

	utils.WriteJSON(w, http.StatusCreated, session.User)
}

func (c *AuthController) LogIn(w http.ResponseWriter, r *http.Request) {
	var input dto.LogInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	session, err := c.authUseCase.LogIn(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	c.setSession(w, session)

	utils.WriteJSON(w, http.StatusOK, session.User)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.RefreshCookieName)
	if err != nil {
		utils.Unauthorized(w)
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
