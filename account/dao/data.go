package dao

import "time"

type CreateCustReq struct {
	Nama string `json:"nama" validate:"required"`
	Nik  string `json:"nik" validate:"required"`
	NoHp string `json:"no_hp" validate:"required"`
	Pin  string `json:"pin" validate:"required"`
	Email  string `json:"email" validate:"required"`
}
type CreateCustRes struct {
	NoRekening *string `json:"no_rekening"`
}

type SaldoRes struct {
	Saldo *float64 `json:"saldo"`
}

type CreateTabungTarikReq struct {
	Nominal    float64 `json:"nominal" validate:"required,gt=0"`
}
type CreateTabungTarikUpdate struct {
	Nominal    float64
	NoRekening    string
	Pin string
}

type CreateTransferReq struct {
	NoRekeningTujuan string  `json:"no_rekening_tujuan" validate:"required"`
	Nominal          float64 `json:"nominal" validate:"required,gt=0"`
}
type CreateTransferUpdate struct {
	NoRekeningAsal string  
	NoRekeningTujuan string  
	Nominal          float64 
}

type PubStruct struct {
	TanggalTransaksi time.Time `json:"tanggal_transaksi"`
	NoRekening string    `json:"no_rekening"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	Nominal    float64   `json:"nominal"`
}

type NoRekeningReq struct {
	NoRekening string  `json:"no_rekening" validate:"required"`
}

type MutasiRes struct {
	Waktu time.Time `json:"waktu"`
	JenisTransaksi string    `json:"kode_transaksi"`
	Nominal float64    `json:"nominal"`
}

type CheckAccountAndPinReq struct {
	NoRekening string
	Pin    string
}

type AccountLoginReq struct {
	NoRekening  string `json:"no_rekening" validate:"required"`
	Pin  string `json:"pin"`
}

type AccountLoginCheckEmailGet struct {
	ID  int
	HashedPin  string
	HashedPassword  string
	NoRekening  string
	NoHp  string
}

type AccountLoginRes struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
}

type JWTField struct{
	NoRekening string
}