package models

import "time"

// Trx represents transaction master table

type Trx struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	IDUser           uint   `gorm:"not null"`
	AlamatPengiriman uint   `gorm:"not null"`
	HargaTotal       int    `gorm:"not null"`
	KodeInvoice      string `gorm:"size:255;unique;not null"`
	MethodBayar      string `gorm:"size:255;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	User      User        `gorm:"foreignKey:IDUser"`
	DetailTrx []DetailTrx `gorm:"foreignKey:IDTrx"`
}
