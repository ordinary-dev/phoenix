package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var Cfg Config

type Config struct {
	// A long and random secret string used for authorization.
	SecretKey string `required:"true"`
	// Path to the sqlite database.
	DBPath string `required:"true"`

	LogLevel string `default:"warning"`

	// Allows you to skip authorization if the "Remote-User" header is specified.
	// Don't use it if you don't know why you need it.
	HeaderAuth bool `default:"false"`

	// Data for the first user.
	// Optional, the site also allows you to create the first user.
	DefaultUsername string
	DefaultPassword string

	// Controls the "secure" option for a token cookie.
	SecureCookie bool `default:"true"`

	// Site title.
	Title string `default:"Phoenix"`
	// Any supported css value, embedded directly into every page.
	FontFamily string `default:"sans-serif"`
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Infof("Config: %v", err)
	}

	err = envconfig.Process("p", &Cfg)
	if err != nil {
		return nil, err
	}

	return &Cfg, nil
}

func (cfg *Config) GetLogLevel() logrus.Level {
	switch cfg.LogLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning", "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.WarnLevel
	}
}
