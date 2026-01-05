package middleware

import (
	"net/http"
	"notes-collab-api/internal/utils"
)

func RequireAuth(jwt *utils.Token) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(utils.AccessCookieName)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userId, err := jwt.ParseAccessToken(cookie.Value)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := utils.WithUserID(r.Context(), userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
