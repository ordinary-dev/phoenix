package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	SecretKey string `required:"true"`
	DBPath    string `required:"true"`
	LogLevel  string `default:"warning"`
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
