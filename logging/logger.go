package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type LoggingSettings struct {
	ServiceName string
	Level       zerolog.Level
}

func InitLogger(settings LoggingSettings) (zerolog.Logger, error) {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Str("service", settings.ServiceName).Logger()
	zerolog.SetGlobalLevel(settings.Level) /// ?
	return logger, nil
}
