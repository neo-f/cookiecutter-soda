package server

import (
	"net/http"
	"time"

	"{{ cookiecutter.project_slug }}"
	"{{ cookiecutter.project_slug }}/internal/config"
	"{{ cookiecutter.project_slug }}/internal/router"
	"{{ cookiecutter.project_slug }}/pkg/tools"
	"{{ cookiecutter.project_slug }}/statics"

	"github.com/neo-f/soda"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

var innerURLs = tools.NewSet(
	config.URL_MONITOR,
	config.URL_HEALTH,
	config.URL_VERSION,
	config.URL_OPENAPI,
	config.URL_REDOC,
	config.URL_SWAGGER,
	config.URL_RAPIDOC,
	config.URL_ELEMENTS,
	config.URL_GATEWAY_CONFIG,
	config.URL_METRICS,
)

var ROUTES = []func(app *soda.Soda){
	RegisterBase,
	router.RegisterDebuggerRouter,
}

// startHttpServer starts configures and starts an HTTP server on the given URL.
// It shuts down the server if any error is received in the error channel.
func InitApp() *soda.Soda {
	app := soda.New("Customer Platform", {{ cookiecutter.project_slug }}.Version,
		soda.EnableValidateRequest(),
		soda.WithOpenAPISpec(config.URL_OPENAPI),
		soda.WithStoplightElements(config.URL_ELEMENTS),
		soda.WithRapiDoc(config.URL_RAPIDOC),
		soda.WithRedoc(config.URL_REDOC),
		soda.WithSwagger(config.URL_SWAGGER),
		soda.WithFiberConfig(
			fiber.Config{
				ReadTimeout:       time.Second * 10,
				EnablePrintRoutes: true,
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					if err == nil {
						return c.Next()
					}
					status := http.StatusInternalServerError //default error status
					if e, ok := err.(*fiber.Error); ok {     // it's a custom error, so use the status in the error
						status = e.Code
					}
					msg := map[string]interface{}{"code": status, "message": err.Error()}
					return c.Status(status).JSON(msg)
				},
			},
		),
	)
	app.OpenAPI().Info.Contact = &openapi3.Contact{Name: "NEO", Email: "tmpgfw@gmail.com"}
	app.OpenAPI().Info.Description = statics.GenDescription()
	for _, env := range config.ENVS {
		app.OpenAPI().AddServer(&openapi3.Server{URL: env[1], Description: env[0]})
	}
	for _, route := range ROUTES {
		route(app)
	}
	return app
}
