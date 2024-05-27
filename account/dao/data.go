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
	Waktu time.Time `json:"waktu"`
	NoRekening string    `json:"no_rekening"`
	IdJurnal          string     `json:"id_jurnal"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	NominalIn      float64      `json:"nominal_in"`
	NominalOut      float64      `json:"nominal_out"`
}

type NoRekeningReq struct {
	NoRekening string  `json:"no_rekening" validate:"required"`
}

type MutasiRes struct {
	Waktu time.Time `json:"waktu"`
	IdJurnal string `json:"id_jurnal"`
	JenisTransaksi string    `json:"kode_transaksi"`
	NominalIn float64    `json:"nominal_in"`
	NominalOut float64    `json:"nominal_out"`
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