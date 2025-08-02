package models

import (
	"time"
)

// User represents the user table
// One-to-many with Alamat, one-to-one with Toko, one-to-many with Trx
// Role flag IsAdmin restricts category management

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Nama         string     `gorm:"size:255;not null"`
	Password     string     `gorm:"size:255;not null"`
	NoTelp       string     `gorm:"size:255;unique;not null"`
	Email        string     `gorm:"size:255;unique;not null"`
	TanggalLahir *time.Time `gorm:"type:date"`
	JenisKelamin string     `gorm:"size:50"`
	Tentang      string     `gorm:"type:text"`
	Pekerjaan    string     `gorm:"size:255"`
	IDProvinsi   string     `gorm:"size:255"`
	IDKota       string     `gorm:"size:255"`
	IsAdmin      bool       `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Alamat      []*Alamat `gorm:"foreignKey:IDUser"`
	Toko        *Toko     `gorm:"foreignKey:IDUser"`
	Trx         []*Trx    `gorm:"foreignKey:IDUser"`
}
