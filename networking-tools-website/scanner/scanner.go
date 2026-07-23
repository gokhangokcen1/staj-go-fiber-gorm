package scanner

import (
	"net"
	"strings"
	"sync"

	"github.com/gokhangokcen1/subnet-backend/oui"
)

type HostSonuc struct {
	IP      string `json:"ip"`
	Ayakta  bool   `json:"ayakta"`
	Mac     string `json:"mac"`
	Uretici string `json:"uretici"`
}

func hostTara(ip string, sonucChannel chan HostSonuc, wg *sync.WaitGroup) {
	defer wg.Done()

	sonuc := HostSonuc{IP: ip}

	mac := MacAdresiBul(ip)

	if ip == myIP() {
		mac = myMAC()
		mac = strings.ReplaceAll(mac, "-", ":")
	}

	ayakta := PingAt(ip)
	if !ayakta && mac != "" {
		ayakta = true
	}

	sonuc.Ayakta = ayakta

	if ayakta {
		sonuc.Mac = mac
		sonuc.Uretici = oui.UreticiBul(mac)
		sonuc.Mac = strings.ReplaceAll(mac, "-", ":")
	}

	sonucChannel <- sonuc
}

func SubnetTara(hostlar []string) []HostSonuc {
	var sonuclar []HostSonuc
	var wg sync.WaitGroup
	sonucChannel := make(chan HostSonuc, len(hostlar))

	for _, ip := range hostlar {
		wg.Add(1)
		go hostTara(ip, sonucChannel, &wg)
	}

	wg.Wait()
	close(sonucChannel)

	for sonuc := range sonucChannel {
		if sonuc.Ayakta {
			sonuclar = append(sonuclar, sonuc)
		}
	}
	return sonuclar
}

func myMAC() string {
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String()
		}
	}
	return ""
}

func myIP() string {
	interfaces, _ := net.Interfaces()

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					return ipv4.String()
				}
			}
		}
	}
	return ""
}
