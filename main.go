package main

import (
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("configuration loading failed", "err", err)
		os.Exit(-1)
	}

	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	})
	slog.SetDefault(slog.New(logger))

	err = database.EstablishDatabaseConnection(cfg)
	if err != nil {
		slog.Error("can't connect to the database", "err", err)
		os.Exit(-1)
	}

	if err := database.ApplyMigrations(); err != nil {
		slog.Error("can't apply database migrations", "err", err)
		os.Exit(-1)
	}

	// Create the first user.
	if cfg.DefaultUsername != "" && cfg.DefaultPassword != "" {
		adminCount, err := database.CountAdmins()
		if err != nil {
			slog.Error("can't query user count", "err", err)
			os.Exit(-1)
		}

		if adminCount < 1 {
			_, err := database.CreateAdmin(cfg.DefaultUsername, cfg.DefaultPassword)
			if err != nil {
				slog.Error("can't create the first user", "err", err)
				os.Exit(-1)
			}
		}
	}

	server, err := web.GetHttpServer()
	if err != nil {
		slog.Error("unable to create a web server", "err", err)
		os.Exit(-1)
	}

	server.ListenAndServe()
}
