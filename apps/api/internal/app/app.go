package app

import (
	"notes-collab-api/internal/config"
	"notes-collab-api/internal/controller"
	"notes-collab-api/internal/repository"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func Run() {
	fx.New(options()).Run()
}

func options() fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger { return errorLogger{} }),
		fx.StopTimeout(10*time.Second),

		config.Module,
		postgresModule,
		repository.Module,
		usecase.Module,
		controller.Module,
		serverModule,

		fx.Provide(
			utils.NewJWT,
			bindSession,
			bindUserRepo,
			bindTokenRepo,
			bindNoteRepo,
			bindWorkspaceRepo,
			bindWorkspaceMemberRepo,
		),

		fx.Invoke(runMigrations),
		fx.Invoke(startServer),
	)
}

func bindSession(t *utils.Token) usecase.Session { return t }

func bindUserRepo(r *repository.UserRepo) usecase.UserRepo { return r }

func bindTokenRepo(r *repository.RefreshTokenRepo) usecase.TokenRepo { return r }

func bindNoteRepo(r *repository.NoteRepo) usecase.NoteRepo { return r }

func bindWorkspaceRepo(r *repository.WorkspaceRepo) usecase.WorkspaceRepo { return r }

func bindWorkspaceMemberRepo(r *repository.WorkspaceMemberRepo) usecase.WorkspaceMemberRepo {
	return r
}
