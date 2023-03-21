package dal

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	SourceField           string
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

// Error implements logger.Interface
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	log.Ctx(ctx).Error().Msgf(s, args...)
}

// Warn implements logger.Interface
func (*GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	log.Ctx(ctx).Warn().Msgf(s, args...)
}

// Info implements logger.Interface
func (*GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	log.Ctx(ctx).Info().Msgf(s, args...)
}

// LogMode implements logger.Interface
func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Trace implements logger.Interface
func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	log := log.Ctx(ctx).
		With().
		Str("sql", sql).
		Str("elapsed", elapsed.String()).
		Int64("rows", rows).
		Logger()

	if l.SourceField != "" {
		log = log.With().Str(l.SourceField, utils.FileWithLineNum()).Logger()
	}
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound) {
		log.Error().Err(err).Msg("[GORM] query error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.Warn().Msg("[GORM] slow query")
		return
	}
	log.Debug().Msg("[GORM] query")
}

var _ logger.Interface = &GormLogger{}

var DefaultGormLogger = &GormLogger{
	SlowThreshold:         time.Second * 3,
	SourceField:           "source",
	SkipErrRecordNotFound: true,
}
