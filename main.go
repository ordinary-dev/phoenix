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
	db, err := database.GetDatabaseConnection(cfg)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	engine := views.GetGinEngine(cfg, db)
	engine.Run(":8080")
}
