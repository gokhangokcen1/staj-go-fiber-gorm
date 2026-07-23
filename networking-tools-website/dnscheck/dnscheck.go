package dnscheck

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"

	"github.com/gokhangokcen1/subnet-backend/models"
)

const queryTimeout = 3 * time.Second

type resolver struct {
	name, address, location string
	latitude, longitude     float64
}

var resolvers = []resolver{
	{"Cloudflare", "1.1.1.1:53", "Global / Cloudflare", 37.7749, -122.4194},
	{"Google Public DNS", "8.8.8.8:53", "Global / Google", 37.4220, -122.0841},
	{"Quad9", "9.9.9.9:53", "Zürich, Switzerland", 47.3769, 8.5417},
	{"OpenDNS", "208.67.222.222:53", "San Francisco, USA", 37.7749, -122.4194},
	{"AdGuard DNS", "94.140.14.14:53", "Frankfurt, Germany", 50.1109, 8.6821},
	{"CleanBrowsing", "185.228.168.9:53", "Los Angeles, USA", 34.0522, -118.2437},
	{"Comodo Secure DNS", "8.26.56.26:53", "Clifton, USA", 40.8584, -74.1638},
	{"Verisign Public DNS", "64.6.64.6:53", "Reston, USA", 38.9586, -77.3570},
	{"Control D", "76.76.2.0:53", "Toronto, Canada", 43.6532, -79.3832},
	{"NextDNS", "45.90.28.0:53", "New York, USA", 40.7128, -74.0060},
	{"DNS.SB", "185.222.222.222:53", "Amsterdam, Netherlands", 52.3676, 4.9041},
	{"Mullvad DNS", "194.242.2.2:53", "Gothenburg, Sweden", 57.7089, 11.9746},
	{"UncensoredDNS", "91.239.100.100:53", "Copenhagen, Denmark", 55.6761, 12.5683},
	{"Alternate DNS", "76.76.19.19:53", "Dallas, USA", 32.7767, -96.7970},
	{"Yandex DNS", "77.88.8.8:53", "Moscow, Russia", 55.7558, 37.6173},
}

var recordTypes = map[string]uint16{
	"A": dns.TypeA, "AAAA": dns.TypeAAAA, "CNAME": dns.TypeCNAME,
	"MX": dns.TypeMX, "NS": dns.TypeNS, "PTR": dns.TypePTR,
	"SRV": dns.TypeSRV, "SOA": dns.TypeSOA, "TXT": dns.TypeTXT,
	"CAA": dns.TypeCAA, "DS": dns.TypeDS, "DNSKEY": dns.TypeDNSKEY,
}

func CheckPropagation(host, requestedType string) (models.DNSPropagationResponse, error) {
	domain := strings.TrimSpace(host)
	if domain == "" {
		return models.DNSPropagationResponse{}, fmt.Errorf("domain zorunludur")
	}
	recordType := strings.ToUpper(strings.TrimSpace(requestedType))
	if recordType == "" {
		recordType = "A"
	}
	dnsType, ok := recordTypes[recordType]
	if !ok {
		return models.DNSPropagationResponse{}, fmt.Errorf("desteklenmeyen kayit turu: %s", recordType)
	}

	response := models.DNSPropagationResponse{
		Domain: domain, RecordType: recordType,
		Results: make([]models.DNSPropagationResult, len(resolvers)),
	}
	var wg sync.WaitGroup
	for index, server := range resolvers {
		wg.Add(1)
		go func(index int, server resolver) {
			defer wg.Done()
			response.Results[index] = queryResolver(domain, dnsType, server)
		}(index, server)
	}
	wg.Wait()
	return response, nil
}

func queryResolver(domain string, recordType uint16, server resolver) models.DNSPropagationResult {
	result := models.DNSPropagationResult{
		Server: server.name, Location: server.location,
		Latitude: server.latitude, Longitude: server.longitude,
	}
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), recordType)
	client := &dns.Client{Timeout: queryTimeout}
	resp, _, err := client.Exchange(msg, server.address)
	if err != nil {
		result.Status, result.Error = "error", "Sunucuya ulasilamadi veya zaman asimina ugradi"
		return result
	}
	if resp == nil || resp.Rcode == dns.RcodeNameError || (resp.Rcode == dns.RcodeSuccess && len(resp.Answer) == 0) {
		result.Status = "not_found"
		return result
	}
	if resp.Rcode != dns.RcodeSuccess {
		result.Status = "error"
		result.Error = "DNS sunucusu hata dondurdu: " + dns.RcodeToString[resp.Rcode]
		return result
	}
	for _, answer := range resp.Answer {
		result.Answers = append(result.Answers, answer.String())
	}
	result.Status = "found"
	return result
}
