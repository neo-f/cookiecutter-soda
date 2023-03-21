package server

import (
	"{{ cookiecutter.project_slug }}"
	"{{ cookiecutter.project_slug }}/internal/config"
	"{{ cookiecutter.project_slug }}/pkg/logger"
	"{{ cookiecutter.project_slug }}/pkg/middlewares/metrics"
	"{{ cookiecutter.project_slug }}/pkg/middlewares/uniform"
	"runtime"
	"time"

	"github.com/neo-f/soda"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type SystemHealth struct {
	MySQL    bool        `json:"mysql"     oai:"description=mysql 连接状态"`
	Redis    bool        `json:"redis"     oai:"description=redis 连接状态"`
	AntsPool interface{} `json:"ants_pool" oai:"description=协程池状态"`
}

func HealthProbe(c *fiber.Ctx) error {
	return c.Status(200).SendString("OK")
}

type SystemVersion struct {
	Version   string `json:"version"    oai:"description=版本"`
	BuildTime string `json:"build_time" oai:"description=构建时间"`
	GoVersion string `json:"go_version" oai:"description=go 版本"`
}

func version(c *fiber.Ctx) error {
	return c.JSON(SystemVersion{
		Version:   {{ cookiecutter.project_slug }}.Version,
		BuildTime: {{ cookiecutter.project_slug }}.BuildTime,
		GoVersion: {{ cookiecutter.project_slug }}.GoVersion,
	})
}

func defaultNextFunc(c *fiber.Ctx) bool {
	return innerURLs.Has(c.Path())
}

func RegisterBase(app *soda.Soda) {
	app.Use(
		// CORS 跨域
		cors.New(),
		// Compress 压缩
		compress.New(),
		// recover 异常恢复
		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Ctx(c.UserContext()).Error().Interface("err", e).Msg(string(buf))
			},
		}),
		// uniform 统一返回结构
		uniform.New(uniform.Config{Next: defaultNextFunc}),
		// requestid 请求增加trace_id
		requestid.New(
			requestid.Config{
				Generator:  func() string { return xid.New().String() },
				ContextKey: "trace_id",
			},
		),
		// kylin ball 麒麟球实现
		logger.NewKylinMiddleware(),
		// logger 请求日志
		flogger.New(flogger.Config{
			Next:       func(c *fiber.Ctx) bool { return innerURLs.Has(c.Path()) },
			TimeFormat: time.RFC3339,
			Format:     "${time} [${locals:trace_id}] ${status} - ${latency} ${method} ${path}\n",
		}),
		// Prometheus
		metrics.NewMiddleware(metrics.Config{Next: defaultNextFunc}),
	)

	if config.Get().Debug {
		app.Use(pprof.New())
		app.Get(config.URL_MONITOR, monitor.New()).
			SetSummary("系统监视器").
			AddTags("System").OK()
	}
	app.Get(config.URL_HEALTH, HealthProbe).
		SetSummary("健康检查").
		SetDescription("\n - 一切正常的时候返回状态码200\n - 任意一项不成功的时候返回状态码400").
		AddJSONResponse(200, SystemHealth{}).
		AddJSONResponse(400, SystemHealth{}).
		AddTags("System").OK()

	app.Get(config.URL_VERSION, version).
		SetSummary("版本信息").
		AddJSONResponse(200, SystemVersion{}).
		AddTags("System").OK()
}
