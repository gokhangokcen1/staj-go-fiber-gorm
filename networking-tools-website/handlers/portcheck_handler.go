package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/gokhangokcen1/subnet-backend/portcheck"
)

func KontrolEtHandler(c fiber.Ctx) error {
	req := new(models.PortCheckRequest)

	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Istek govdesi okunamadi",
		})
	}

	if req.IP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "IP adresi zorunludur",
		})
	}

	if req.Port < 1 || req.Port > 65535 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Port 1 ile 65535 arasinda olmalidir",
		})
	}

	acik, detay := portcheck.KontrolEt(req.IP, req.Port)

	response := models.PortCheckResponse{
		IP:    req.IP,
		Port:  req.Port,
		Acik:  acik,
		Detay: detay,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
