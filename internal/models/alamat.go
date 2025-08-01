package models

import "time"

// Alamat represents user shipping addresses

type Alamat struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	IDUser       uint   `gorm:"not null"`
	JudulAlamat  string `gorm:"size:255"`
	NamaPenerima string `gorm:"size:255"`
	NoTelp       string `gorm:"size:255"`
	DetailAlamat string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	User *User `gorm:"foreignKey:IDUser"`
}
