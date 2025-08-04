package models

import "time"

// DetailTrx represents transaction detail rows

type DetailTrx struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	IDTrx       uint `gorm:"not null"`
	IDLogProduk uint `gorm:"not null"`
	IDToko      uint `gorm:"not null"`
	Kuantitas   int  `gorm:"not null"`
	HargaTotal  int  `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
