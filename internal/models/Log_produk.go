// File: internal/models/log_produk.go

package models

import "time"

// LogProduk snapshots product data at a specific moment (e.g. creation or checkout)
type LogProduk struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	IDProduk      uint      `gorm:"not null"`           // FK → Produk
	NamaProduk    string    `gorm:"size:255;not null"`  // salin dari produk master
	Slug          string    `gorm:"size:255;not null"`  // salin dari produk master
	HargaReseller string    `gorm:"size:255;not null"`  // harga reseller pada saat snapshot
	HargaKonsumen string    `gorm:"size:255;not null"`  // harga konsumen pada saat snapshot
	Deskripsi     string    `gorm:"type:text;not null"` // salin dari produk master
	IDToko        uint      `gorm:"not null"`           // FK → Toko
	IDCategory    uint      `gorm:"not null"`           // salin dari produk master
	StokAwal      int       `gorm:"not null"`           // stok pada saat snapshot
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	// Relasi ke Produk master (opsional untuk preload)
	Produk Produk `gorm:"foreignKey:IDProduk;references:ID"`
}
