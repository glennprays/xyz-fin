package model

import "time"

type Transaction struct {
	NomorKontrak  string    `json:"nomor_kontrak"`
	ConsumerNIK   string    `json:"consumer_nik"`
	OTR           float64   `json:"otr"`
	AdminFee      float64   `json:"admin_fee"`
	JumlahCicilan int       `json:"jumlah_cicilan"`
	JumlahBunga   float64   `json:"jumlah_bunga"`
	NamaAsset     string    `json:"nama_asset"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransactionRequest struct {
	ConsumerNIK string  `json:"consumer_nik"`
	OTR         float64 `json:"otr"`
	Tenor       int     `json:"tenor"`
	NamaAsset   string  `json:"nama_asset"`
}

type TransactionResponse struct {
	NomorKontrak  string  `json:"nomor_kontrak"`
	ConsumerNIK   string  `json:"consumer_nik"`
	OTR           float64 `json:"otr"`
	AdminFee      float64 `json:"admin_fee"`
	JumlahCicilan int     `json:"jumlah_cicilan"`
	JumlahBunga   float64 `json:"jumlah_bunga"`
	NamaAsset     string  `json:"nama_asset"`
	Status        string  `json:"status"`
}
