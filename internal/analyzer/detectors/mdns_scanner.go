package detectors

import (
	"hattic/internal/analyzer"
	"strings"
	"sync"
	"time"
)

type MDNSScanner struct {
	*analyzer.BaseDetector
	knownIPs map[string]bool
	mu       sync.Mutex
}

func NewMDNSScanner() *MDNSScanner {
	return &MDNSScanner{
		BaseDetector: analyzer.NewBaseDetector(5 * time.Minute),
		knownIPs:     make(map[string]bool),
	}
}

func (d *MDNSScanner) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" || strings.ToUpper(packet.Protocol) != "MDNS" {
		return alerts
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.knownIPs[packet.SrcIP] {
		d.knownIPs[packet.SrcIP] = true
		hostname := packet.Hostname
		if hostname == "" {
			hostname = "Unknown Hostname"
		}
		alerts = append(alerts, analyzer.Alert{
			Timestamp: packet.Timestamp.Format(time.RFC3339),
			Level:     analyzer.LevelInfo,
			Type:      "NEW_DEVICE_MDNS",
			Message:   "New device discovered via mDNS: " + packet.SrcIP + " (" + hostname + ")",
			SrcIP:     packet.SrcIP,
		})
	}

	return alerts
}
