package dao

import "time"

type SubStruct struct {
	TanggalTransaksi time.Time `json:"tanggal_transaksi"`
	NoRekening string    `json:"no_rekening"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	Nominal    float64   `json:"nominal"`
}
type CreateJournalRes struct {
	Success bool `json:"success"`
}
