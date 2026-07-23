package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	listenAddr := ":" + port

	app := fiber.New()

	app.Use(cors.New())

	routes.SetupRoutes(app)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint bulunamadi",
		})
	})

	log.Printf("Sunucu %s adresinde başlatılıyor...\n", listenAddr)
	if err := app.Listen(listenAddr); err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Fatalf("Port %s zaten kullanımda. Bu portu serbest bırakın veya PORT değişkenini ayarlayarak farklı bir port kullanın. Örnek: $env:PORT='3002'; go run main.go", port)
		}
		log.Fatal(err)
	}
}
