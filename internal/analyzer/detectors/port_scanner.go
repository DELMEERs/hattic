package detectors

import (
	"fmt"
	"hattic/internal/analyzer"
	"sync"
	"time"
)

type PortScannerDetector struct {
	*analyzer.BaseDetector
	portThreshold int
	ipToPorts     map[string]map[uint16]bool
	mu            sync.Mutex
	whitelistIPs  map[string]bool
	ignorePorts   map[uint16]bool
}

func NewPortScannerDetector(threshold int) *PortScannerDetector {
	return &PortScannerDetector{
		BaseDetector:  analyzer.NewBaseDetector(5 * time.Minute),
		portThreshold: threshold,
		ipToPorts:     make(map[string]map[uint16]bool),
		whitelistIPs: map[string]bool{
			"192.168.0.1": true,
			"127.0.0.1":   true,
		},
		ignorePorts: map[uint16]bool{
			1900: true,
			5353: true,
			3702: true,
		},
	}
}

func (d *PortScannerDetector) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" || packet.DstPort == 0 {
		return alerts
	}

	if d.whitelistIPs[packet.SrcIP] || d.ignorePorts[packet.DstPort] {
		return alerts
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.ipToPorts[packet.SrcIP]; !exists {
		d.ipToPorts[packet.SrcIP] = make(map[uint16]bool)
	}

	d.ipToPorts[packet.SrcIP][packet.DstPort] = true
	portCount := len(d.ipToPorts[packet.SrcIP])

	if portCount >= d.portThreshold {
		if d.ShouldAlert("PORT_SCAN", packet.SrcIP) {
			alerts = append(alerts, analyzer.Alert{
				Timestamp: packet.Timestamp.Format(time.RFC3339),
				Level:     analyzer.LevelWarning,
				Type:      "PORT_SCAN",
				Message:   fmt.Sprintf("IP %s is scanning ports: %d unique ports hit.", packet.SrcIP, portCount),
				SrcIP:     packet.SrcIP,
			})
		}
	}

	return alerts
}
