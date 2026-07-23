package oui

import (
	"encoding/csv"
	"os"
	"strings"
)

var OuiHaritasi = make(map[string]string)

func Yukle(dosyaYolu string) error {
	dosya, err := os.Open(dosyaYolu)
	if err != nil {
		return err
	}

	defer dosya.Close()

	okuyucu := csv.NewReader(dosya)
	okuyucu.LazyQuotes = true
	okuyucu.FieldsPerRecord = -1

	kayitlar, err := okuyucu.ReadAll()
	if err != nil {
		return err
	}

	for i, kayit := range kayitlar {
		if i == 0 || len(kayit) < 3 {
			continue
		}
		assignment := strings.ToUpper(strings.TrimSpace(kayit[1]))
		organizasyon := strings.TrimSpace(kayit[2])
		if assignment != "" && organizasyon != "" {
			OuiHaritasi[assignment] = organizasyon
		}
	}

	return nil
}

func UreticiBul(mac string) string {
	if mac == "" {
		return "Bilinmiyor"
	}

	temiz := strings.ReplaceAll(mac, "-", "")
	temiz = strings.ReplaceAll(temiz, ":", "")
	temiz = strings.ToUpper(temiz)

	if len(temiz) < 6 {
		return "Bilinmiyor"
	}

	oui := temiz[:6]
	if uretici, bulundu := OuiHaritasi[oui]; bulundu {
		return uretici
	}
	return "Bilinmiyor"
}
