package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {

	env := strings.ToLower(os.Getenv("APP_ENV"))
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// DEV
	if env == "dev" {
		writer := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC822,
		}
		Log = zerolog.New(writer).With().Timestamp().Logger()
		return
	}

	// TEST
	if env == "test" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
		return
	}

	// PROD — JSON
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
