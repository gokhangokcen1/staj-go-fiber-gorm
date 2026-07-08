package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/subnet", handlers.HesaplaSubnet)
}
