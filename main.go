package main

import (
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	logLevel := cfg.GetLogLevel()
	log.SetLevel(logLevel)
	log.Infof("Setting log level to %v", logLevel)

	err = database.EstablishDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.ApplyMigrations(); err != nil {
		log.Fatal(err)
	}

	// Create the first user.
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

	server, err := web.GetHttpServer()
	if err != nil {
		log.Fatal(err)
	}

	server.ListenAndServe()
}
