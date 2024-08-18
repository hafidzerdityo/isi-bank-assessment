package dao

import "time"

type AccountWithCustomerQuery struct {
	// Fields from Account
	ID          int       `json:"id"`
	IdNasabah   int       `json:"id_nasabah"`
	NoRekening  string    `json:"no_rekening"`
	HashedPin   string    `json:"hashed_pin"`
	Saldo       float64   `json:"saldo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	IsDeleted   bool      `json:"is_deleted"`
	// Fields from Customer
	CustomerID    int       `json:"customer_id"`
	Nama          string    `json:"nama"`
	Nik           string    `json:"nik"`
	NoHp          string    `json:"no_hp"`
	Email         string    `json:"email"`
	CustomerCreatedAt time.Time `json:"customer_created_at"`
	CustomerUpdatedAt *time.Time `json:"customer_updated_at"`
	CustomerIsDeleted bool `json:"customer_is_deleted"`
}

type CreateTabungTarikUpdate struct {
	Nominal    float64
	NoRekening    string
}

type CreateTransferUpdate struct {
	NoRekeningAsal string  
	NoRekeningTujuan string  
	Nominal          float64 
}
