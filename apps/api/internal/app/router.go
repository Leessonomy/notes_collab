package app

import (
	"net/http"
	"notes-collab-api/internal/config"
	"notes-collab-api/internal/controller"
	"notes-collab-api/internal/middleware"
	"notes-collab-api/internal/utils"
)

func buildHandlers(
	cfg *config.Config,
	jwt *utils.Token,
	authCtrl *controller.AuthController,
	noteCtrl *controller.NoteController,
	workspaceCtrl *controller.WorkspaceController,
) http.Handler {
	protected := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(jwt)(h)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handleHealth)

	mux.HandleFunc("POST /api/auth/signup", authCtrl.SignUp)
	mux.HandleFunc("POST /api/auth/login", authCtrl.LogIn)
	mux.HandleFunc("POST /api/auth/refresh", authCtrl.Refresh)
	mux.HandleFunc("POST /api/auth/logout", authCtrl.LogOut)

	mux.Handle("GET /api/auth/me", protected(authCtrl.GetMe))
	mux.Handle("POST /api/auth/change-password", protected(authCtrl.ChangePassword))

	mux.Handle("GET /api/notes/{id}", protected(noteCtrl.GetByID))
	mux.Handle("POST /api/notes", protected(noteCtrl.CreateNote))
	mux.Handle("PATCH /api/notes/{id}", protected(noteCtrl.UpdateNote))
	mux.Handle("DELETE /api/notes/{id}", protected(noteCtrl.DeleteNote))

	mux.Handle("GET /api/workspaces", protected(workspaceCtrl.ListWorkspacesWithNotes))
	mux.Handle("POST /api/workspaces", protected(workspaceCtrl.CreateWorkspace))
	mux.Handle("PATCH /api/workspaces/{id}", protected(workspaceCtrl.UpdateWorkspace))
	mux.Handle("DELETE /api/workspaces/{id}", protected(workspaceCtrl.DeleteWorkspace))

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
