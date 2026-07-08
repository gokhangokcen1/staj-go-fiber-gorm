package database

import (
	"fmt"
	"log"

	"github.com/gokhangokcen1/ogrenci-crud/config"
	"github.com/gokhangokcen1/ogrenci-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// her yerden erisilebilir gorm bağlantısı > handlers'ta da kullanılacak
var DB *gorm.DB

// ConnectDB, config'ten bilgileri alır PostgreSQL'e bağlanır
// modeldeki tablo yapısını veritabanına uygular (automigrate)
func ConnectDB(cfg config.AppConfig) {
	// dsn: data source name
	// postgresql'e nasıl bağlanacağız -> değerler
	// kendi bilgisayarımız oldugu icin sslmode=disable, > require
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	// gorm.Open(postgres.Open(dsn)) : verdiğimiz bilgilerle dsn'e bağlan
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Logger: her sorguyu komut satırına yazdır.
	})

	if err != nil {
		log.Fatal("Veritabanina baglanilamadi: ", err)
	}

	log.Println("PostgreSQL baglantisi basarili")
	log.Println("Migration calistiriliyor...")

	//AutoMigrate: struct'ı inceler, yoksa oluşturur, güncelleme varsa ekler.
	err = db.AutoMigrate(&models.Ogrenci{})
	if err != nil {
		log.Fatal("Migration basarisiz: ", err)
	}
	log.Println("Migration tamamlandi")
	DB = db
}
