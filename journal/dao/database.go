package dao

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(user string, 
	dbname string, 
	password string, 
	host string,
	port string,) (*gorm.DB, error) {
	// Replace with your PostgreSQL connection details
	dsn := "user=" + user + " dbname=" + dbname + " password=" + password + " host=" + host + " port="+ port + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate your models here
	err = db.AutoMigrate(&Journal{})
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
