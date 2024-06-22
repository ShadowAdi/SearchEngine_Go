package routes

import (
	"fmt"
	"search_engine/db"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

func SetRoutes(app *fiber.App) {
	fmt.Println("Connectoion strarted")
	app.Get("/", AuthMiddlewar, DashboardHandler)
	app.Post("/", AuthMiddlewar, DashboardPostHandler)
	app.Get("/Login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Get("/User", func(c *fiber.Ctx) error {
		u := &db.User{}
		u.CreateAdmin()
		return c.SendString("User Created")
	})

}
