package dao

import "time"

type PubStruct struct {
	Waktu time.Time `json:"waktu"`
	NoRekening string    `json:"no_rekening"`
	IdJurnal          string     `json:"id_jurnal"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	NominalIn      float64      `json:"nominal_in"`
	NominalOut      float64      `json:"nominal_out"`
}