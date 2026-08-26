package app

import (
	"net/http"
	"notes-collab-api/internal/controller"
	"notes-collab-api/internal/middleware"
	"notes-collab-api/internal/repository"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createHandler(cfg *config, pool *pgxpool.Pool) http.Handler {
	jwt := utils.NewJWT(cfg.JWTSecret, cfg.AccessExpire, cfg.RefreshExpire)

	userRepo := repository.NewUserRepo(pool)
	refreshRepo := repository.NewRefreshTokenRepo(pool)
	workspaceRepo := repository.NewWorkspaceRepo(pool)
	noteRepo := repository.NewNoteRepo(pool)

	authUC := usecase.NewAuth(userRepo, refreshRepo, jwt)
	workspaceUC := usecase.NewWorkspace(workspaceRepo, noteRepo)
	noteUC := usecase.NewNote(noteRepo)

	authCtrl := controller.NewAuthController(authUC, cfg.CookieSecure)
	workspaceCtrl := controller.NewWorkspaceController(workspaceUC)
	noteCtrl := controller.NewNoteController(noteUC)

	protectedAuth := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(jwt)(h)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handleHealth)

	mux.HandleFunc("POST /api/auth/signup", authCtrl.SignUp)
	mux.HandleFunc("POST /api/auth/login", authCtrl.LogIn)
	mux.HandleFunc("POST /api/auth/refresh", authCtrl.Refresh)
	mux.HandleFunc("POST /api/auth/logout", authCtrl.LogOut)

	mux.Handle("GET /api/auth/me", protectedAuth(authCtrl.GetMe))

	mux.Handle("POST /api/notes", protectedAuth(noteCtrl.CreateNote))
	mux.Handle("GET /api/notes/{id}", protectedAuth(noteCtrl.GetByID))
	mux.Handle("DELETE /api/notes/{id}", protectedAuth(noteCtrl.DeleteNote))

	mux.Handle("GET /api/workspaces", protectedAuth(workspaceCtrl.ListWorkspacesWithNotes))
	mux.Handle("POST /api/workspaces", protectedAuth(workspaceCtrl.CreateWorkspace))
	mux.Handle("DELETE /api/workspaces/{id}", protectedAuth(workspaceCtrl.DeleteWorkspace))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", cfg.CORSOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Add("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		mux.ServeHTTP(w, r)
	})
}
