package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/gokhangokcen1/subnet-backend/sslcheck"
)

func SslKontrol(c fiber.Ctx) error {
	req := new(models.SslCheckRequest)

	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "İstek gövdesi okunamadı.",
		})
	}

	req.Website = strings.TrimSpace(req.Website)

	if req.Website == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Website adresi zorunludur.",
		})
	}

	if req.Port != 443 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SSL kontrolü yalnızca 443 portu için yapılabilir.",
		})
	}

	report, err := sslcheck.CheckSSL(req.Website)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(report)
}
