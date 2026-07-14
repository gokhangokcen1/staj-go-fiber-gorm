package handlers

import (
	"context"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/capture"
	"github.com/gokhangokcen1/subnet-backend/models"
)

var (
	currentPacketChan chan models.PacketInfo
	currentCancel     context.CancelFunc
)

func StartCaptureHandler(c fiber.Ctx) error {
	req := new(models.StartCaptureRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Istek govdesi okunamadi"})
	}
	if req.Device == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cihaz adi zorunludur"})
	}

	if currentCancel != nil {
		currentCancel()
	}

	ch, cancel, err := capture.StartCapture(req.Device)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Capture baslatilamadi", "detay": err.Error()})
	}

	currentPacketChan = ch
	currentCancel = cancel
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"mesaj": "Capture baslatildi"})
}

func StopCaptureHandler(c fiber.Ctx) error {
	if currentCancel != nil {
		currentCancel()
		currentCancel = nil
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"mesaj": "Capture durduruldu"})
}

func ListDevicesHandler(c fiber.Ctx) error {
	devices, err := capture.ListDevices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cihazlar alinamadi"})
	}

	var result []models.DeviceInfo
	for _, d := range devices {
		result = append(result, models.DeviceInfo{Name: d.Name, Description: d.Description})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func CaptureWebSocketHandler(c *websocket.Conn) {
	for {
		if currentPacketChan == nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		info, ok := <-currentPacketChan
		if !ok {
			return
		}
		if err := c.WriteJSON(info); err != nil {
			return
		}
	}
}
