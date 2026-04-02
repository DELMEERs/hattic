package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hattic/internal/db"
	"hattic/internal/sniffer"
)

func main() {
	device := flag.String("device", "eth0", "Network interface to sniff on")
	dbPath := flag.String("db", "data/traffic.db", "Path to the SQLite database")
	flag.Parse()

	// Initialize Database
	database, err := db.InitDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Channel for Log Entries
	logChan := make(chan db.TrafficLog, 100)

	// Start Database Writer Goroutine
	go func() {
		for logEntry := range logChan {
			if err := database.Create(&logEntry).Error; err != nil {
				log.Printf("Error saving log to database: %v", err)
			}
		}
	}()

	// Handle Signals for Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start Sniffer
	s := sniffer.NewSniffer(*device, logChan)
	log.Printf("Starting sniffer on device: %s", *device)

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Sniffer failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down engine...")
	close(logChan)
}
