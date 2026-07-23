package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/gokhangokcen1/subnet-backend/whois"
)

func WhoisLookupHandler(c fiber.Ctx) error {
	req := new(models.WhoisRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Istek govdesi okunamadi"})
	}
	result, err := whois.Lookup(strings.TrimSpace(req.Domain))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}
