package main

import (
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/views"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configure logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Read config
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	// Set log level
	logLevel := cfg.GetLogLevel()
	logrus.SetLevel(logLevel)
	logrus.Infof("Setting log level to %v", logLevel)

	// Connect to the database
	err = database.EstablishDatabaseConnection(cfg)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	// Create the first user
	if cfg.DefaultUsername != "" && cfg.DefaultPassword != "" && database.CountAdmins() < 1 {
		_, err := database.CreateAdmin(cfg.DefaultUsername, cfg.DefaultPassword)
		if err != nil {
			logrus.Errorf("%v", err)
		}
	}

	server, err := views.GetHttpServer()
	if err != nil {
		logrus.Fatal(err)
	}

	server.ListenAndServe()
}
