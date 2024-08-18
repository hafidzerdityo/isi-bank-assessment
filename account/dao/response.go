package dao

import "time"

type CreateCustRes struct {
	NoRekening *string `json:"no_rekening"`
}

type SaldoRes struct {
	Saldo *float64 `json:"saldo"`
}

type SaldoWithCifRes struct {
	Saldo *float64 `json:"saldo"`
	Nama  string   `json:"nama"`
}

type MutasiData struct {
	Waktu          time.Time `json:"waktu"`
	IdJurnal       string    `json:"id_jurnal"`
	JenisTransaksi string    `json:"kode_transaksi"`
	NominalIn      float64   `json:"nominal_in"`
	NominalOut     float64   `json:"nominal_out"`
}

type MutasiRes struct {
	ListData  []MutasiData `json:"list_data"`
	TotalData int64        `json:"total_data"`
}

type AccountLoginRes struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
}


type JWTField struct{
	NoRekening string
}
