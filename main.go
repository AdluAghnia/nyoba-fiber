package main

import (
	"log"

	"github.com/AdluAghnia/nyoba-fiber/handler"
	middlewares "github.com/AdluAghnia/nyoba-fiber/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/auth/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Login": "logim",
		}, "layouts/main")
	})
	// create a new JWT middleware
	// Routes
	app.Get("/", handler.FrontpageHandler)
	app.Post("/login", handler.LoginHandler)

	app.Use(middlewares.JwtMiddleware)
	app.Get("/protected", handler.Protected)

	// app.Use("/content", jwtMiddleware, func(c *fiber.Ctx) error {
	// 	return c.Render("front", fiber.Map{
	// 		"Title": "This IS CONTENT !!!",
	// 	}, "layouts/main")
	// })

	log.Fatal(app.Listen(":6969"))
}
