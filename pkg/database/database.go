package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GlobalBan struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Reason    string
}
type GlobalBans []GlobalBan

type DB struct {
	Id      int
	Name    string
	Content string
}

var dsn = "host=192.168.3.42 port=5432 user=gobot_canary password=dev dbname=gobot_canary sslmode=disable TimeZone=Asia/Tokyo"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default})

func GetDBConn() (*gorm.DB, error) {
	if err != nil {
		return db, err
	}
	return db, nil
}
