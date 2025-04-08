package api

import (
	"fmt"

	"github.com/MagicRodri/go_graphql_service/internal/db"
	"github.com/MagicRodri/go_graphql_service/internal/logging"
	"github.com/gofiber/fiber/v2"
)

func SetupServer(host string, port int) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Powered-By", "go-graphql")
		return c.Next()
	})

	// Setup routes
	registerHandlers(app)

	// launch server
	if port > 0 {
		port = 8080
	}
	if host == "" {
		host = "localhost"
	}
	address := fmt.Sprintf("%s:%d", host, port)
	logging.Get().Infof("Server starting on %s", address)
	logging.Get().Fatal(app.Listen(address))
}

func registerHandlers(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		db.GetDB().Ping()
		return c.SendStatus(fiber.StatusOK)
	})
	app.Post("/query", graphqlHandler)
}
