package models

import "time"

type Category struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	NamaCategory string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Produk []Produk `gorm:"foreignKey:IDCategory"`
}
