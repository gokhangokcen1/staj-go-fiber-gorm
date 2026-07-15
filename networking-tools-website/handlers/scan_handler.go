package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/gokhangokcen1/subnet-backend/scanner"
	"github.com/gokhangokcen1/subnet-backend/subnet"
)

func TaraHandler(c fiber.Ctx) error {
	req := new(models.ScanRequest)

	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Istek govdesi okunamadi"})
	}
	if req.IP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IP adresi zorunludur"})
	}
	if req.CIDR < 1 || req.CIDR > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CIDR 1 ile 30 arasinda olmalidir"})
	}

	hostlar := subnet.HostIPListesi(req.IP, req.CIDR)

	if len(hostlar) > 1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bu CIDR cok fazla host icerir, daha dar bir subnet secin (ornegin /22 veya daha kucuk aralik)",
		})
	}

	sonuclar := scanner.SubnetTara(hostlar)

	return c.Status(fiber.StatusOK).JSON(sonuclar)
}
