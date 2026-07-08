package models

import "gorm.io/gorm"

// Ogrenci, veritabanındaki "ogrencis" tablsouna karsilik gelen struct

type Ogrenci struct {
	gorm.Model
	Ad     string `json:"ad"`
	Soyad  string `json:"soyad"`
	Numara string `json:"numara" gorm:"unique"`
	Bolum  string `json:"bolum"`
}

// Otomatik olarak 4 alan eklenir:
// - ID
// - CreatedAt : kayıt ne zaman oluşturuldu
// - UpdatedAt : kayıt en son ne zaman güncellendi
// - DeletedAt : kayıt silindi mi
