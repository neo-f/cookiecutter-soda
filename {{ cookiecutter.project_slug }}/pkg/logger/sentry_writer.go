package logger

import (
	"io"
	"time"

	"{{ cookiecutter.project_slug }}/internal/config"

	// "github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

var levelsMapping = map[zerolog.Level]sentry.Level{
	zerolog.DebugLevel: sentry.LevelDebug,
	zerolog.InfoLevel:  sentry.LevelInfo,
	zerolog.WarnLevel:  sentry.LevelWarning,
	zerolog.ErrorLevel: sentry.LevelError,
	zerolog.FatalLevel: sentry.LevelFatal,
	zerolog.PanicLevel: sentry.LevelFatal,
}

var _ = io.WriteCloser(new(SentryWriter))

var now = time.Now

// SentryWriter is a sentry events writer with std io.SentryWriter iface.
type SentryWriter struct {
	client *sentry.Client

	levels       map[zerolog.Level]struct{}
	flushTimeout time.Duration
}

// Write handles zerolog's json and sends events to sentry.
func (w *SentryWriter) Write(data []byte) (int, error) {
	event, ok := w.parseLogEvent(data)
	if ok {
		w.client.CaptureEvent(event, nil, nil)
		// should flush before os.Exit
		if event.Level == sentry.LevelFatal {
			w.client.Flush(w.flushTimeout)
		}
	}

	return len(data), nil
}

// Close forces client to flush all pending events.
// Can be useful before application exits.
func (w *SentryWriter) Close() error {
	w.client.Flush(w.flushTimeout)
	return nil
}

func (w *SentryWriter) parseLogEvent(data []byte) (*sentry.Event, bool) {
	const logger = "zerolog"
	lvlStr := gjson.GetBytes(data, zerolog.LevelFieldName)
	lvl, err := zerolog.ParseLevel(lvlStr.String())
	if err != nil {
		return nil, false
	}

	_, enabled := w.levels[lvl]
	if !enabled {
		return nil, false
	}

	sentryLvl, ok := levelsMapping[lvl]
	if !ok {
		return nil, false
	}

	event := sentry.Event{
		Timestamp: now(),
		Level:     sentryLvl,
		Logger:    logger,
		Tags:      make(map[string]string, 6),
		Request:   &sentry.Request{},
	}

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		switch key.String() {
		// case zerolog.LevelFieldName, zerolog.TimestampFieldName:
		case zerolog.MessageFieldName:
			event.Message = value.String()
		case zerolog.ErrorFieldName:
			event.Exception = append(event.Exception, sentry.Exception{
				Value:      value.String(),
				Stacktrace: newStacktrace(),
			})
		case config.LogTagURL:
			event.Request.URL = value.String()
		case config.LogTagMethod:
			event.Request.Method = value.String()
		case config.LogTagHeaders:
			headers := make(map[string]string)
			value.ForEach(func(key, value gjson.Result) bool {
				headers[key.String()] = value.String()
				return true
			})
			event.Request.Headers = headers
		case config.LogTagData:
			event.Request.Data = value.String()
		case config.LogTagAuthorization:
			event.Tags["Authorization"] = value.String()
		case config.LogTagBID:
			event.Tags["bid"] = value.String()
		case config.LogTagStaffID:
			event.Tags["staff_id"] = value.String()
		case config.LogTagTraceID:
			event.Tags["trace_id"] = value.String()
		default:
			event.Tags[key.String()] = value.String()
		}
		return true
	})

	return &event, true
}

func newStacktrace() *sentry.Stacktrace {
	const (
		module       = "github.com/archdx/zerolog-sentry"
		loggerModule = "github.com/rs/zerolog"
	)

	st := sentry.NewStacktrace()

	threshold := len(st.Frames) - 1
	// drop current module frames
	for ; threshold > 0 && st.Frames[threshold].Module == module; threshold-- {
	}

outer:
	// try to drop zerolog module frames after logger call point
	for i := threshold; i > 0; i-- {
		if st.Frames[i].Module == loggerModule {
			for j := i - 1; j >= 0; j-- {
				if st.Frames[j].Module != loggerModule {
					threshold = j
					break outer
				}
			}

			break
		}
	}

	st.Frames = st.Frames[:threshold+1]

	return st
}

// WriterOption configures sentry events writer.
type WriterOption interface {
	apply(*options)
}

type optionFunc func(*options)

func (fn optionFunc) apply(c *options) { fn(c) }

type options struct {
	release      string
	environment  string
	serverName   string
	levels       []zerolog.Level
	ignoreErrors []string
	sampleRate   float64
	flushTimeout time.Duration
	debug        bool
}

// WithLevels configures zerolog levels that have to be sent to Sentry.
// Default levels are: error, fatal, panic.
func WithLevels(levels ...zerolog.Level) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.levels = levels
	})
}

// WithSampleRate configures the sample rate as a percentage of events to be sent in the range of 0.0 to 1.0.
func WithSampleRate(rate float64) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.sampleRate = rate
	})
}

// WithRelease configures the release to be sent with events.
func WithRelease(release string) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.release = release
	})
}

// WithEnvironment configures the environment to be sent with events.
func WithEnvironment(environment string) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.environment = environment
	})
}

// WithServerName configures the server name field for events. Default value is OS hostname.
func WithServerName(serverName string) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.serverName = serverName
	})
}

// WithIgnoreErrors configures the list of regexp strings that will be used to match against event's message
// and if applicable, caught errors type and value. If the match is found, then a whole event will be dropped.
func WithIgnoreErrors(reList []string) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.ignoreErrors = reList
	})
}

// WithDebug enables sentry client debug logs.
func WithDebug(debug bool) WriterOption {
	return optionFunc(func(cfg *options) {
		cfg.debug = debug
	})
}

// NewSentryWriter creates writer with provided DSN and options.
func NewSentryWriter(dsn string, opts ...WriterOption) (*SentryWriter, error) {
	cfg := newDefaultConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:          dsn,
		SampleRate:   cfg.sampleRate,
		Release:      cfg.release,
		Environment:  cfg.environment,
		ServerName:   cfg.serverName,
		IgnoreErrors: cfg.ignoreErrors,
		Debug:        cfg.debug,
	})
	if err != nil {
		return nil, err
	}

	levels := make(map[zerolog.Level]struct{}, len(cfg.levels))
	for _, lvl := range cfg.levels {
		levels[lvl] = struct{}{}
	}

	return &SentryWriter{
		client:       client,
		levels:       levels,
		flushTimeout: cfg.flushTimeout,
	}, nil
}

func newDefaultConfig() options {
	return options{
		levels: []zerolog.Level{
			zerolog.ErrorLevel,
			zerolog.FatalLevel,
			zerolog.PanicLevel,
		},
		sampleRate:   1.0,
		flushTimeout: 3 * time.Second,
	}
}
