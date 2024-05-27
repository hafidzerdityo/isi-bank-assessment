package dao

import (
	"time"
)


type Journal struct{
	ID          int       `gorm:"primaryKey"`
	NoRekening          string    `gorm:"size:255;"`
	IdJurnal          string       `gorm:"not null"`
	JenisTransaksi          string    `gorm:"size:1;not null"`
	NominalIn      *float64     `gorm:"type:numeric(10,2)"`
	NominalOut      *float64     `gorm:"type:numeric(10,2)"`
	Waktu     time.Time 
}


func (Journal) TableName() string {
    return "journal"
}


