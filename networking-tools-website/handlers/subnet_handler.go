package handlers

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/gokhangokcen1/subnet-backend/subnet"
)

func HesaplaSubnet(c fiber.Ctx) error {
	req := new(models.SubnetRequest)

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

	if req.CIDR < 0 || req.CIDR > 32 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "CIDR 0 ile 32 arasinda olmalidir",
		})
	}
	if parsedIP := net.ParseIP(req.IP); parsedIP == nil || parsedIP.To4() == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz IPv4 adresi",
		})
	}

	ipBolunmus := strings.Split(req.IP, ".")
	if len(ipBolunmus) != 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz IP formati (Örnek: 192.168.1.1)",
		})
	}

	binaryIP := subnet.DecToBinary(ipBolunmus)
	subnetMask := subnet.FindSubnetBinary(req.CIDR)
	wildcardMask := subnet.WildcardMask(subnetMask)
	netAddress := subnet.NetworkAddress(binaryIP, req.CIDR)
	broadAddress := subnet.BroadcastAddress(binaryIP, req.CIDR)
	ilkKullanilabilir := subnet.IlkKullanilabilir(netAddress)
	sonKullanilabilir := subnet.SonKullanilabilir(broadAddress)
	hostSayisi := subnet.HostSayisi(req.CIDR)

	response := models.SubnetResponse{
		Adres:            req.IP,
		NetworkMask:      subnet.BinaryToIP(subnetMask),
		WildcardMask:     subnet.BinaryToIP(wildcardMask),
		NetworkAddress:   subnet.BinaryToIP(netAddress),
		BroadcastAddress: subnet.BinaryToIP(broadAddress),
		Hostmin:          subnet.BinaryToIP(ilkKullanilabilir),
		Hostmax:          subnet.BinaryToIP(sonKullanilabilir),
		HostsPerNet:      hostSayisi,

		AddressBinary:      subnet.BinaryToString(binaryIP, req.CIDR),
		MaskBinary:         subnet.BinaryToString(subnetMask, req.CIDR),
		WildcardMaskBinary: subnet.BinaryToString(wildcardMask, req.CIDR),
		NetworkBinary:      subnet.BinaryToString(netAddress, req.CIDR),
		BroadcastBinary:    subnet.BinaryToString(broadAddress, req.CIDR),
		HostMinBinary:      subnet.BinaryToString(ilkKullanilabilir, req.CIDR),
		HostMaxBinary:      subnet.BinaryToString(sonKullanilabilir, req.CIDR),
		Class:              subnet.GetIPClass(req.IP),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func getLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}

			if ip[0] == 169 && ip[1] == 254 {
				continue
			}

			return ip.String()
		}
	}

	return ""
}

func MevcutIP(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"ip": getLocalIP(),
	})
}
