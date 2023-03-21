package uniform

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// 调整JSON结果返回为公司要求的统一格式

// New creates a new middleware handler
func New(cfgs ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(cfgs...)

	// Return new handler
	return func(ctx *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}

		uniformedResp := []byte{}
		_ = ctx.Next()
		code := ctx.Response().StatusCode()
		if code >= 200 && code < 300 {
			if string(ctx.Response().Header.ContentType()) == "application/json" {
				uniformedResp, _ = sjson.SetBytes(uniformedResp, "code", 0)
				uniformedResp, _ = sjson.SetRawBytes(uniformedResp, "data", ctx.Response().Body())
				ctx.Response().SetBodyRaw(uniformedResp)
			}
			return nil
		}
		if code == 422 {
			uniformedResp, _ = sjson.SetBytes(uniformedResp, "code", 422)
			uniformedResp, _ = sjson.SetBytes(uniformedResp, "message", "参数校验失败")
			uniformedResp, _ = sjson.SetRawBytes(uniformedResp, "details", ctx.Response().Body())
			ctx.Response().SetBodyRaw(uniformedResp)
			return nil
		}
		if code >= 400 {
			log.Ctx(ctx.UserContext()).Error().
				Int("status", code).
				Str("request", ctx.String()).
				Bytes("body", uniformedResp).
				Msg("error occurred")
			uniformedResp, _ = sjson.SetBytes(uniformedResp, "code", code)

			resp := ctx.Context().Response.Body()
			if msg := gjson.GetBytes(resp, "message"); msg.Exists() {
				uniformedResp, _ = sjson.SetBytes(uniformedResp, "message", msg.String())
			} else {
				uniformedResp, _ = sjson.SetBytes(uniformedResp, "message", resp)
			}
			ctx.Response().SetBodyRaw(uniformedResp)
		}
		return nil
	}
}
