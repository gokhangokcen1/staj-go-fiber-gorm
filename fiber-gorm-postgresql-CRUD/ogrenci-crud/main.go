package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gokhangokcen1/ogrenci-crud/config"
	"github.com/gokhangokcen1/ogrenci-crud/database"
	"github.com/gokhangokcen1/ogrenci-crud/routes"
)

func main() {
	cfg := config.LoadConfig()

	database.ConnectDB(cfg)

	app := fiber.New()

	app.Use(cors.New())

	routes.SetupRoutes(app)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint bulunamadi",
		})
	})

	log.Fatal(app.Listen(":" + cfg.AppPort))

}
