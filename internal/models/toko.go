package models

import "time"

// Toko represents the toko table, automatically created upon user registration

type Toko struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	IDUser    uint   `gorm:"not null;unique"`
	NamaToko  string `gorm:"size:255"`
	URLFoto   string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User   User     `gorm:"foreignKey:IDUser"`
	Produk []Produk `gorm:"foreignKey:IDToko"`
}
