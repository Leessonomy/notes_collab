package server

import (
	"context"
	"net/http"
	"notes-collab-api/internal/config"
	"notes-collab-api/internal/controller"
	"notes-collab-api/internal/db"
	"notes-collab-api/internal/middleware"
	"notes-collab-api/internal/repository"
	"notes-collab-api/internal/utils"
)

func New(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	pool, err := db.CreatePostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(ctx, pool); err != nil {
		return nil, err
	}

	jwt := utils.NewJWT(cfg.JWTSecret, cfg.AccessExpire, cfg.RefreshExpire)

	userRepo := repository.NewUserRepo(pool)
	refreshRepo := repository.NewRefreshTokenRepo(pool)
	workspaceRepo := repository.NewWorkspaceRepo(pool)
	noteRepo := repository.NewNoteRepo(pool)

	authCtrl := controller.NewAuthController(userRepo, refreshRepo, jwt, cfg.CookieSecure)
	workspaceCtrl := controller.NewWorkspaceController(workspaceRepo)
	noteCtrl := controller.NewNoteController(noteRepo)

	requireAuth := middleware.RequireAuth(jwt)
	protected := func(h http.HandlerFunc) http.Handler { return requireAuth(h) }

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handleHealth)

	mux.HandleFunc("POST /api/auth/signup", authCtrl.SignUp)
	mux.HandleFunc("POST /api/auth/login", authCtrl.LogIn)
	mux.HandleFunc("POST /api/auth/refresh", authCtrl.Refresh)
	mux.HandleFunc("POST /api/auth/logout", authCtrl.LogOut)

	mux.Handle("POST /api/notes", protected(noteCtrl.Create))

	mux.Handle("GET /api/workspaces", protected(workspaceCtrl.GetAll))
	mux.Handle("POST /api/workspaces", protected(workspaceCtrl.Create))
	mux.Handle("GET /api/workspaces/{workspaceId}/notes", protected(noteCtrl.GetByWorkspaceID))

	return withCORS(cfg.CORSOrigin, mux), nil
}

func withCORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Add("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
