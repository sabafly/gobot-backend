package database

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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

type TransMCServer struct {
	FeedMCServer
	Address string
	Port    uint16
}

type FeedMCServer struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Hash      string
	GuildID   string
	ChannelID string
	RoleID    string
	Name      string
	Locale    discordgo.Locale
	PanelID   string `gorm:"primarykey"`
}

type FeedMCServers []FeedMCServer

type DB struct {
	Id      int
	Name    string
	Content string
}

var dsn string
var db *gorm.DB
var err error

func init() {
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	dsn = "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable TimeZone=Asia/Tokyo"
	log.Print(dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default})
}

func GetDBConn() (*gorm.DB, error) {
	if err != nil {
		return db, err
	}
	return db, nil
}
