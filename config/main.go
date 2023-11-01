package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	SecretKey       string `required:"true"`
	DBPath          string `required:"true"`
	LogLevel        string `default:"warning"`
	EnableGinLogger bool   `default:"false"`
	Production      bool   `default:"true"`
	HeaderAuth      bool   `default:"false"`
	DefaultUsername string
	DefaultPassword string
	// Controls the "secure" option for a token cookie.
	SecureCookie bool `default:"true"`
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Infof("Config: %v", err)
	}

	var cfg Config
	err = envconfig.Process("p", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) GetLogLevel() logrus.Level {
	switch cfg.LogLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.WarnLevel
	}
}
