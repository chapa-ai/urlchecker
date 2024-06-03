package app

import (
	"github.com/chapa-ai/urlchecker/config"
	"github.com/chapa-ai/urlchecker/internal/checks"
	"github.com/chapa-ai/urlchecker/logging"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

type App interface {
	Serve()
}

type app struct {
	log zerolog.Logger
	cfg config.Config
}

func (a *app) Serve() {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	go func(logger zerolog.Logger, cfg config.Config) {
		err := checks.DoChecksWithInterval(logger, a.cfg)
		if err != nil {
			a.log.Fatal().Err(err).Str("failed", "DoChecksWithInterval").Msgf("failed DoChecksWithInterval: %v", err)
		}
	}(a.log, a.cfg)

	a.log.Info().Msg("app started successfully")

	sig := <-shutdownSignal
	a.log.Info().Msgf("Received signal...: %v", sig)
}

func NewApp(cfg config.Config) App {
	loggingSettings := logging.LoggingSettings{
		ServiceName: cfg.Logging.LogServiceName,
		Level:       cfg.Logging.LogLevel,
	}

	logger, err := logging.InitLogger(loggingSettings)
	if err != nil {
		panic(err)
	}

	return &app{
		log: logger,
		cfg: cfg,
	}
}
