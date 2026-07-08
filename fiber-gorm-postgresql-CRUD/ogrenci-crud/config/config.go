package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// her seferinde db'e bağlanmak için bilgileri yazmaktansa bunu otomatize hale getiriyoruz ve db'den bilgileri otomatik olarak .env üzerinden alıyoruz. Eğer ki .env'de bir problem çıkarsa da ortam değişkenlerinden bu bilgileri çekiyoruz.

type AppConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Uyari: .env dosyasi bulunamadi ortam degiskenlerinden okunacak")
	}

	return AppConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		AppPort:    os.Getenv("APP_PORT"),
	}
}
