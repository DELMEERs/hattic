package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TrafficLog struct {
	ID          uint      `gorm:"primaryKey"`
	Timestamp   time.Time `gorm:"index"`
	SrcIP       string
	DstIP       string
	SrcMAC      string
	DstMAC      string
	Protocol    string
	PayloadSize int
	PacketCount int
}

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&TrafficLog{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
