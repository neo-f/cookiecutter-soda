package logger

import (
	"context"
	"strings"
	"time"

	"{{ cookiecutter.project_slug }}/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
  CtxDebug    = struct{}{}
  CtxTimezone = struct{}{}
)

func NewKylinMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctxLogger := log.With().Interface(config.LogTagTraceID, c.Locals("trace_id"))

		if authorization := c.Get(fiber.HeaderAuthorization); authorization != "" {
			ctxLogger.Str(config.LogTagAuthorization, authorization)
		}
		if bid := c.Get("Rcrai-Bid"); bid != "" {
			ctxLogger.Str(config.LogTagBID, bid)
		}
		if staffID := c.Get("Rcrai-StaffId"); staffID != "" {
			ctxLogger.Str(config.LogTagStaffID, staffID)
		}
		ctx := c.UserContext()
		logger := ctxLogger.Logger()
		logger = logger.Level(zerolog.InfoLevel)
		if lv := c.Get("Debug_Kylin"); strings.ToLower(lv) == "true" {
			ctx = context.WithValue(ctx, CtxDebug, struct{}{})
			logger = logger.Level(zerolog.DebugLevel)
		}
		if tz := c.Get("R-Timezone"); tz != "" {
			if loc, err := time.LoadLocation(tz); err == nil {
				ctx = context.WithValue(ctx, CtxTimezone, loc)
			}
		}
		c.SetUserContext(logger.WithContext(ctx))
		return c.Next()
	}
}
