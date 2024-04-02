package main

import (
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/views"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Read config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set log level
	logLevel := cfg.GetLogLevel()
	log.SetLevel(logLevel)
	log.Infof("Setting log level to %v", logLevel)

	// Connect to the database
	err = database.EstablishDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Apply migrations.
	if err := database.ApplyMigrations(); err != nil {
		log.Fatal(err)
	}

	// Create the first user
	if cfg.DefaultUsername != "" && cfg.DefaultPassword != "" {
		adminCount, err := database.CountAdmins()
		if err != nil {
			log.Fatal(err)
		}

		if adminCount < 1 {
			_, err := database.CreateAdmin(cfg.DefaultUsername, cfg.DefaultPassword)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	server, err := views.GetHttpServer()
	if err != nil {
		log.Fatal(err)
	}

	server.ListenAndServe()
}
