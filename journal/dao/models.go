package dao

import (
	"time"
)


type Journal struct{
	ID          int       `gorm:"primaryKey"`
	TanggalTransaksi     time.Time 
	NoRekening          string    `gorm:"size:255;"`
	Nominal      float64     `gorm:"type:numeric(10,2);not null"`
	JenisTransaksi          string    `gorm:"size:10;"`
}


func (Journal) TableName() string {
    return "journal"
}


