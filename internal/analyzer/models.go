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

type PacketInfo struct {
	Timestamp time.Time `json:"timestamp"`
	SrcIP     string    `json:"src_ip"`
	DstIP     string    `json:"dst_ip"`
	SrcMAC    string    `json:"src_mac"`
	DstMAC    string    `json:"dst_mac"`
	SrcPort   uint16    `json:"src_port"`
	DstPort   uint16    `json:"dst_port"`
	Protocol  string    `json:"protocol"`
	TTL       uint8     `json:"ttl"`
	Hostname  string    `json:"hostname"`
	Length    int       `json:"length"`
}
