# Kodu çalıştırma
```
go mod tidy
go run main.go
```

# Test komutları 
## Öğrenci oluşturma 
`curl -X POST http://localhost:3000/api/ogrenci -H "Content-Type: application/json" -d "{\"ad\":\"Ahmet\",\"soyad\":\"Yilmaz\",\"numara\":\"2021001\",\"bolum\":\"Bilgisayar Muh\"}"`

## Tüm öğrencileri listele 
`curl http://localhost:3000/api/ogrenci`

## Tek öğrenci getir 
`curl http://localhost:3000/api/ogrenci/1`

## Öğrenci güncelle 
`curl -X PUT http://localhost:3000/api/ogrenci/1 -H "Content-Type: application/json" -d "{\"bolum\":\"Elektrik Elektronik Muh\"}"`

## Öğrenci sil 
`curl -X DELETE http://localhost:3000/api/ogrenci/1`  