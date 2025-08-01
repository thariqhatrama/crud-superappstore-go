package models

import "time"

// LogProduk snapshots product data per transaction

type LogProduk struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	IDProduk      uint   `gorm:"not null"`
	NamaProduk    string `gorm:"size:255;not null"`
	Slug          string `gorm:"size:255;not null"`
	HargaReseller string `gorm:"size:255;not null"`
	HargaKonsumen string `gorm:"size:255;not null"`
	Deskripsi     string `gorm:"type:text"`
	IDToko        uint   `gorm:"not null"`
	IDCategory    uint   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Produk Produk `gorm:"foreignKey:IDProduk"`
}
