package main

import (
	"fiberioc/configs"
	"fiberioc/routes"

	_ "fiberioc/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	// app.Server().MaxConnsPerIP = 1
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	routes.Router(app)

	var port string = configs.GetPort()
	app.Listen(":" + port)
}
