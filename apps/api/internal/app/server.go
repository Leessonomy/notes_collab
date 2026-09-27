package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"notes-collab-api/internal/config"
	"time"

	"go.uber.org/fx"
)

var serverModule = fx.Module("server",
	fx.Provide(
		buildHandlers,
		newServer,
	),
)

func newServer(cfg *config.Config, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func startServer(lc fx.Lifecycle, srv *http.Server, sd fx.Shutdowner) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}

			log.Printf("listening on %s", srv.Addr)

			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("server error: %v", err)
					_ = sd.Shutdown(fx.ExitCode(1))
				}
			}()

			return nil
		},
		OnStop: srv.Shutdown,
	})
}
