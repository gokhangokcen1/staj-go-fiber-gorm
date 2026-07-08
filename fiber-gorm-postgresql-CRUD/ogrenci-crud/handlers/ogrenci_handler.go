package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gokhangokcen1/ogrenci-crud/database"
	"github.com/gokhangokcen1/ogrenci-crud/models"
)

// CreateOgrenci: yeni ogrenci kaydi
func CreateOgrenci(c fiber.Ctx) error {
	ogrenci := new(models.Ogrenci)

	// bind, body: gelen HTTP isteğin body'sindeki JSON'ı ogrenci struct'ına doldurur.
	if err := c.Bind().Body(ogrenci); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Istek govdesi okunamadi",
		})
	}

	if ogrenci.Ad == "" || ogrenci.Numara == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ad ve Numara alanlari zorunludur",
		})
	}

	result := database.DB.Create(&ogrenci) // GORM -> INSERT
	// oluşturma başarılı olursa ID ve CreatedAt struct'a yazılır
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ogrenci olusturulamadi",
			"detay": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ogrenci)
	// 201 Created koduyla oluşturulan kayıt json olarak döner
}

// GetAllOgrenciler: tum ogrencileri listele
func GetAllOgrenciler(c fiber.Ctx) error {
	var ogrenciler []models.Ogrenci

	// SELECT * FROM ogrencis WHERE deleted_at IS NULL
	result := database.DB.Find(&ogrenciler)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ogrenciler getirilemedi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ogrenciler)

}

// GetOgrenci: id'ye göre bir öğrenci getir
func GetOgrenci(c fiber.Ctx) error {

	// URL'den gelen :id paramatresi stringtir. Biz bunu int'e çeviriz.
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz id",
		})
	}

	var ogrenci models.Ogrenci

	// WHERE id = ? LIMIT 1
	result := database.DB.First(&ogrenci, id)

	if result.Error != nil {
		// 404 Not Found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ogrenci bulunamadi",
		})

	}
	return c.Status(fiber.StatusOK).JSON(ogrenci)
}

// UpdateOgrenci: var olan ogrenciyi guncelle
func UpdateOgrenci(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz id",
		})
	}

	// Oncelikle ogrenci var mi diye kontrol ediyoruz
	var ogrenci models.Ogrenci
	if result := database.DB.First(&ogrenci, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ogrenci bulunamadi",
		})
	}

	// direkt orijinal struct üzerinde guncelleme yapmaktansa öncelikle ayrı bir guncelleme struct'ı oluşturuyoruz.
	// Cunku, bind().body() her alanı doldurmazsa mevcut degerleri sifirlayabilir
	guncelleme := new(models.Ogrenci)

	// http req. body'sindeki json bilgilerini tekrardan struct'a yazar
	if err := c.Bind().Body(guncelleme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Istek govdesi okunamadi",
		})
	}

	// Updates metodu, guncelleme struct'ındaki sıfır olmayan alanları gunceller
	result := database.DB.Model(&ogrenci).Updates(guncelleme)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Guncelleme basarisiz",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ogrenci)
}

// DeleteOgrenci: bir ogrenciyi siler (soft delete)
func DeleteOgrenci(c fiber.Ctx) error {
	// Yine id string'ten int'e çevrildi.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz id",
		})
	}

	// soft delete uygulanır ve DeletedAt alanı doldurulur
	// veriler fiziksel olarak silinmez yalnızca deletedat kısmı boş olan veriler sorgularda gözükür
	result := database.DB.Unscoped().Delete(&models.Ogrenci{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Silme basarisiz",
		})
	}

	// olmayan bir id'yi silemeyeceğimizden, ogrencinin olduğunu kontrol ediyoruz
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ogrenci bulunamadi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"mesaj": "Ogrenci silindi",
	})
}
