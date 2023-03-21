package logger

import (
	"io"
	"os"

	"{{ cookiecutter.project_slug }}"
	"{{ cookiecutter.project_slug }}/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Setup() {
	cfg := config.Get()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	writers := []io.Writer{}
	if config.Get().IsLocal {
		writers = append(writers, zerolog.NewConsoleWriter())
	} else {
		writers = append(writers, os.Stderr)
	}
	if dsn := cfg.Sentry.DSN; dsn != "" {
		sentryWriter, err := NewSentryWriter(
			cfg.Sentry.DSN,
			WithDebug(cfg.Debug),
			WithEnvironment(cfg.Sentry.Environment),
			WithRelease({{ cookiecutter.project_slug }}.Version),
			WithServerName("{{ cookiecutter.project_slug }}-service"),
		)
		if err == nil {
			writers = append(writers, sentryWriter)
		}
	}
	log.Logger = zerolog.New(io.MultiWriter(writers...)).With().Stack().Timestamp().Logger()
}
