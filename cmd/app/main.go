package main

import (
	"FinalTask/config"
	"FinalTask/internal/models"
	"FinalTask/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inisialisasi DB & JWT secret
	config.InitDB()
	config.InitSecret()

	// AutoMigrate semua tabel
	config.DB.AutoMigrate(
		&models.User{},
		&models.Alamat{},
		&models.Toko{},
		&models.Category{},
		&models.Produk{},
		&models.FotoProduk{},
		&models.LogProduk{},
		&models.Trx{},
		&models.DetailTrx{},
	)

	app := fiber.New()

	// Routes
	router.SetupRoutes(app)

	app.Listen(":8000")
}
