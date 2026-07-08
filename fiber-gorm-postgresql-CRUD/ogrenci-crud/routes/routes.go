package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/ogrenci-crud/handlers"
)

// SetupRoutes, tum endpoint tanimlarini tek bir yerde toplar.
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/ogrenci", handlers.CreateOgrenci)
	api.Get("/ogrenci", handlers.GetAllOgrenciler)
	api.Get("/ogrenci/:id", handlers.GetOgrenci)
	api.Put("/ogrenci/:id", handlers.UpdateOgrenci)
	api.Delete("/ogrenci/:id", handlers.DeleteOgrenci)
}
