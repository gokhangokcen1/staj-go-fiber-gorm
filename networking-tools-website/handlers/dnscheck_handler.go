package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/gokhangokcen1/subnet-backend/dnscheck"
	"github.com/gokhangokcen1/subnet-backend/models"
)

func CheckDNSHandler(c fiber.Ctx) error {
	req := new(models.DNSCheckRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Istek govdesi okunamadi"})
	}
	if req.Domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Domain zorunludur"})
	}

	result, err := dnscheck.CheckPropagation(req.Domain, req.RecordType)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
