package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Level  string
	Pretty bool
}

var (
	Logger zerolog.Logger //nolint:gochecknoglobals // global by design
)

func New(cfg *Config) zerolog.Logger {
	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	return zerolog.New(output).With().Timestamp().Logger()
}

func InitGlobalLogger(cfg *Config) {
	Logger = New(cfg)
	Logger.Info().Str("logger_level", Logger.GetLevel().String()).Msg("logger initialized")
}
