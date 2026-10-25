package app

import (
	"net/http"
	"notes-collab-api/internal/config"
	"notes-collab-api/internal/controller"
	"notes-collab-api/internal/middleware"
	"notes-collab-api/internal/token"
)

func newHandler(
	cfg *config.Config,
	tokens *token.Manager,
	authCtrl *controller.AuthController,
	noteCtrl *controller.NoteController,
	workspaceCtrl *controller.WorkspaceController,
) http.Handler {
	protected := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(tokens)(h)
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

	return middleware.CORS(cfg.CORSOrigin)(mux)
}
