package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Start Database Writer Goroutine (Batch Insert)
	go func() {
		var buffer []db.TrafficLog
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		flush := func() {
			if len(buffer) > 0 {
				if err := db.BatchCreate(database, buffer); err != nil {
					log.Printf("Error during batch insert: %v", err)
				}
				buffer = nil
			}
		}

		for {
			select {
			case logEntry, ok := <-logChan:
				if !ok {
					flush()
					return
				}
				buffer = append(buffer, logEntry)
				if len(buffer) >= 100 {
					flush()
				}
			case <-ticker.C:
				flush()
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
