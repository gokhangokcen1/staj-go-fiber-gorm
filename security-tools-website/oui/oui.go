package oui

import (
	"encoding/csv"
	"os"
	"strings"
)

var OuiHaritasi = make(map[string]string)

func Yukle(dosyaYolu string) error {
	// fmt.Println("YUKLE ÇALIŞTI")
	dosya, err := os.Open(dosyaYolu)
	if err != nil {
		return err
	}

	defer dosya.Close()

	okuyucu := csv.NewReader(dosya)
	// okuyucu.Comma = '\t'
	okuyucu.LazyQuotes = true
	okuyucu.FieldsPerRecord = -1

	kayitlar, err := okuyucu.ReadAll()
	// fmt.Println("Toplam satır:", len(kayitlar))

	// for i := 0; i < 5 && i < len(kayitlar); i++ {
	// 	fmt.Printf("SATIR %d: %#v\n", i, kayitlar[i])
	// }
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

	// for i, kayit := range kayitlar {

	// 	if i < 5 {
	// 		fmt.Printf("Satir %d: %#v\n", i, kayit)
	// 	}

	// 	if i == 0 || len(kayit) < 3 {
	// 		continue
	// 	}

	// 	assignment := strings.ToUpper(strings.TrimSpace(kayit[1]))
	// 	organizasyon := strings.TrimSpace(kayit[2])

	// 	if assignment != "" && organizasyon != "" {
	// 		OuiHaritasi[assignment] = organizasyon

	// 		if len(OuiHaritasi) <= 5 {
	// 			fmt.Println("Eklendi:", assignment, "->", organizasyon)
	// 		}
	// 	}
	// }

	// fmt.Println("Toplam OUI:", len(OuiHaritasi))
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
	// fmt.Println(OuiHaritasi)
	if uretici, bulundu := OuiHaritasi[oui]; bulundu {
		// fmt.Println("uretici: ", uretici)
		// fmt.Println("bulundu: ", bulundu)
		return uretici
	}
	return "Bilinmiyor"
}
