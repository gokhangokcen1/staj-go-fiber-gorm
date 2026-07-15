package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gokhangokcen1/subnet-backend/oui"
	"github.com/gokhangokcen1/subnet-backend/routes"
)

func main() {
	err := oui.Yukle("data/oui.csv")
	if err != nil {
		log.Fatal("OUI dosyasi yuklenemedi: ", err)
	}
	log.Printf("OUI veritabani yuklendi: %d kayit\n", len(oui.OuiHaritasi))
	fmt.Println("Toplam kayıt:", len(oui.OuiHaritasi))

	app := fiber.New()

	app.Use(cors.New())

	routes.SetupRoutes(app)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint bulunamadi",
		})
	})

	log.Fatal(app.Listen(":3001"))
}
