package models

import "time"

// FotoProduk stores image URLs for products

type FotoProduk struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	IDProduk  uint   `gorm:"not null"`
	URL       string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Produk Produk `gorm:"foreignKey:IDProduk"`
}
