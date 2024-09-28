package dao

import "time"

type Journal struct {
    NoRekening    string      `gorm:"primaryKey;size:255;"`
    IdJurnal      string      `gorm:"primaryKey;not null;size:255"`
    JenisTransaksi string     `gorm:"size:1;not null"`
    NominalIn     *float64    `gorm:"type:numeric(20,2)"`
    NominalOut    *float64    `gorm:"type:numeric(20,2)"`
    Waktu         time.Time 
}

func (Journal) TableName() string {
    return "journal"
}


