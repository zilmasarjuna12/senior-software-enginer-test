package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the logger with proper configuration
func InitLogger() {
	Log = logrus.New()

	// Set log level based on environment
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		Log.SetLevel(logrus.DebugLevel)
	case "INFO":
		Log.SetLevel(logrus.InfoLevel)
	case "WARN":
		Log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if os.Getenv("ENV") == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	Log.SetOutput(os.Stdout)
}

// WithRequestID creates a logger entry
func WithRequestID(ctx context.Context) *logrus.Entry {
	requestID, _ := ctx.Value("request_id").(string)
	return Log.WithField("request_id", requestID)
}

// WithFields creates a logger entry with custom fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithError creates a logger entry with error
func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}

// GetLogger returns the main logger instance
func GetLogger() *logrus.Logger {
	return Log
}
