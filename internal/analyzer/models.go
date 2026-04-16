package analyzer

import "time"

type Level string

const (
	LevelInfo     Level = "Info"
	LevelWarning  Level = "Warning"
	LevelCritical Level = "Critical"
)

type Alert struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	SrcIP     string `json:"src_ip"`
}

// represents extracted information from a network packet
type PacketInfo struct {
	Timestamp time.Time
	SrcIP     string
	DstIP     string
	SrcMAC    string
	DstMAC    string
	SrcPort   uint16
	DstPort   uint16
	Protocol  string
	TTL       uint8
	Hostname  string // for mDNS
	Length    int
}
