package main

import (
	"context"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func handleInterrupt(srv *http.Server, connsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	// We received an interrupt signal, shut down.
	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Info("HTTP server shutdown", "err", err)
	}

	close(connsClosed)
}

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
		userCount, err := database.CountUsers()
		if err != nil {
			slog.Error("can't query user count", "err", err)
			os.Exit(-1)
		}

		if userCount < 1 {
			_, err := database.CreateUser(cfg.DefaultUsername, &cfg.DefaultPassword)
			if err != nil {
				slog.Error("can't create the first user", "err", err)
				os.Exit(-1)
			}
		}
	}

	handler, err := web.GetHandler()
	if err != nil {
		slog.Error("unable to create a web server", "err", err)
		os.Exit(-1)
	}

	http.Handle("/", handler)

	var listener net.Listener
	if cfg.SocketPath != "" {
		slog.Info("starting a web server", "address", cfg.SocketPath)
		listener, err = net.Listen("unix", cfg.SocketPath)
	} else {
		slog.Info("starting a web server", "address", cfg.ListeningAddress)
		listener, err = net.Listen("tcp", cfg.ListeningAddress)
	}

	if err != nil {
		slog.Error("unable to start a web server", "err", err)
		os.Exit(-1)
	}

	connsClosed := make(chan struct{})

	var srv http.Server
	go handleInterrupt(&srv, connsClosed)

	if err = srv.Serve(listener); err != http.ErrServerClosed {
		slog.Error("http server returned an error", "err", err)
		os.Exit(-1)
	}

	<-connsClosed
}
