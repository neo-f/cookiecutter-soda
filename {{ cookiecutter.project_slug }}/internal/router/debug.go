package router

import (
	"context"
	"{{ cookiecutter.project_slug }}/statics"
	"fmt"
	"sort"

	"github.com/neo-f/soda"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func RegisterDebuggerRouter(app *soda.Soda) {
	app.App.Get("/debugger", basicAuth, func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(statics.DebuggerHTML)
	})
	app.App.Get("/debugger/functions", basicAuth, func(c *fiber.Ctx) error {
		sort.Slice(debugFunctions, func(i, j int) bool {
			return debugFunctions[i].Title < debugFunctions[j].Title
		})
		return c.JSON(debugFunctions)
	})
	app.App.Post("/debugger/execute", basicAuth, func(c *fiber.Ctx) error {
		var params struct {
			Func   string            `json:"func"`
			Params map[string]string `json:"params"`
		}
		if err := c.BodyParser(&params); err != nil {
			return err
		}
		var fn *DebugFunction
		for _, f := range debugFunctions {
			if f.Name == params.Func {
				fn = &f
				break
			}
		}
		if fn == nil {
			return fmt.Errorf("function not found")
		}
		resp, err := fn.Do(c.UserContext(), params.Params)
		if err != nil {
			return err
		}
		return c.JSON(resp)
	})
}

var basicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{"neo": "whosyourdaddy"},
})

type DebugFunctionParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type DebugFunction struct {
	Name        string                                                        `json:"name"`
	Title       string                                                        `json:"title"`
	Description string                                                        `json:"description"`
	Parameters  []DebugFunctionParameter                                      `json:"parameters"`
	Do          func(context.Context, map[string]string) (interface{}, error) `json:"-"`
}

var debugFunctions = []DebugFunction{
	{
		Name:        "echo",
		Title:       "demo debugger",
		Description: "echo the parameters",
		Parameters: []DebugFunctionParameter{
			{Name: "param-a", Description: "param-a", Required: true},
			{Name: "param-b", Description: "param-b", Required: true},
		},
		Do: func(ctx context.Context, params map[string]string) (interface{}, error) {
			return params, nil
		},
	},
}
