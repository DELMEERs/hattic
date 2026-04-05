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
	SrcPort     int
	DstPort     int
	Protocol    string
	TTL         int
	PayloadSize int
	PacketCount int
}

type Alert struct {
	ID        uint      `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"index"`
	Level     string
	Type      string
	Message   string
	SrcIP     string
}

func InitDB(dbPath string) (*gorm.DB, error) {
	// Enable WAL mode for better concurrency
	dsn := dbPath + "?_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&TrafficLog{}, &Alert{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func BatchCreate(db *gorm.DB, logs []TrafficLog) error {
	if len(logs) == 0 {
		return nil
	}
	return db.Create(&logs).Error
}
