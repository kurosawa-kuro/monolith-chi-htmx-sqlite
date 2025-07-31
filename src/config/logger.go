package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// LoadLoggerConfig loads logger configuration from environment
func LoadLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
		Output: getEnv("LOG_OUTPUT", "stdout"),
	}
}

// SetupLogger configures the global logger
func SetupLogger(cfg *LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	switch cfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output
	switch cfg.Output {
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "stdout":
		logger.SetOutput(os.Stdout)
	default:
		logger.SetOutput(os.Stdout)
	}

	return logger
}
