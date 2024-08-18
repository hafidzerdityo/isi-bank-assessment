package dao

type CreateCustReq struct {
	Nama  string `json:"nama" validate:"required"`
	Nik   string `json:"nik" validate:"required"`
	NoHp  string `json:"no_hp" validate:"required"`
	Pin   string `json:"pin" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type CreateTabungTarikReq struct {
	Nominal float64 `json:"nominal" validate:"required,gt=0"`
}

type CreateTransferReq struct {
	NoRekeningTujuan string  `json:"no_rekening_tujuan" validate:"required"`
	Nominal          float64 `json:"nominal" validate:"required,gt=0"`
}

type NoRekeningReq struct {
	NoRekening string `json:"no_rekening" validate:"required"`
}

type MutasiReq struct {
	NoRekening string `json:"no_rekening" validate:"required"`
	Page       *int   `json:"page"`
	Limit      *int   `json:"limit"`
}

type AccountLoginReq struct {
	NoRekening string `json:"no_rekening" validate:"required"`
	Pin        string `json:"pin"`
}