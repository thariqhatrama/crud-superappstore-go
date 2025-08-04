package models

import "time"

// Trx represents transaction master table

type Trx struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	IDUser           uint      `json:"id_user"`
	AlamatPengiriman uint      `json:"alamat_pengiriman"`
	MethodBayar      string    `json:"method_bayar"`
	HargaTotal       int       `json:"harga_total"`
	KodeInvoice      string    `json:"kode_invoice"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	DetailTrx []DetailTrx `gorm:"foreignKey:IDTrx" json:"detail_trx"`
}
