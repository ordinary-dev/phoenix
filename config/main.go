package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getStrFromEnv(defaultValue string, varNames ...string) string {
	for i, name := range varNames {
		res := os.Getenv(name)
		if res == "" {
			continue
		}

		if i > 0 {
			msg := fmt.Sprintf("'%s' is deprecated, use '%s' instead", varNames[i], varNames[0])
			slog.Warn(msg)
		}

		return res
	}

	return defaultValue
}

func getBoolFromEnv(defaultValue bool, varNames ...string) bool {
	res := getStrFromEnv("", varNames...)
	if res == "" {
		return defaultValue
	}

	res = strings.ToLower(res)
	return res == "true" || res == "1"
}

var Cfg Config

type Config struct {
	// Path to the sqlite database.
	DBPath string

	LogLevel string

	// Allows you to skip authorization if the "Remote-User" header is specified.
	// Don't use it if you don't know why you need it.
	HeaderAuth bool

	// Data for the first user (optional).
	DefaultUsername string
	DefaultPassword string
	// Allow registration via web interface?
	EnableRegistration bool

	// Controls the "secure" option for a token cookie.
	SecureCookie bool

	// Site title.
	Title string
	// Any supported css value, embedded directly into every page.
	FontFamily string
}

func GetConfig() (*Config, error) {
	godotenv.Load()

	Cfg.DBPath = getStrFromEnv("", "DB_PATH", "P_DBPATH")
	if Cfg.DBPath == "" {
		return nil, errors.New("database path is undefined, set it using DB_PATH environment variable")
	}

	Cfg.LogLevel = getStrFromEnv("warning", "LOG_LEVEL", "P_LOGLEVEL")
	Cfg.HeaderAuth = getBoolFromEnv(false, "HEADER_AUTH", "P_HEADERAUTH")

	Cfg.DefaultUsername = getStrFromEnv("", "DEFAULT_USERNAME", "P_DEFAULTUSERNAME")
	Cfg.DefaultPassword = getStrFromEnv("", "DEFAULT_PASSWORD", "P_DEFAULTPASSWORD")
	Cfg.EnableRegistration = getBoolFromEnv(true, "ENABLE_REGISTRATION")

	Cfg.SecureCookie = getBoolFromEnv(true, "SECURE_COOKIE", "P_SECURECOOKIE")
	Cfg.Title = getStrFromEnv("Phoenix", "TITLE", "P_TITLE")
	Cfg.FontFamily = getStrFromEnv("sans-serif", "FONT_FAMILY", "P_FONTFAMILY")

	return &Cfg, nil
}

func (cfg *Config) GetLogLevel() slog.Level {
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warning", "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}
