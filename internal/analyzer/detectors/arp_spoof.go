package detectors

import (
	"hattic/internal/analyzer"
	"strings"
	"sync"
	"time"
)

type ARPSpoofDetector struct {
	*analyzer.BaseDetector
	ipToMacs map[string]map[string]bool
	mu       sync.Mutex
}

func NewARPSpoofDetector() *ARPSpoofDetector {
	return &ARPSpoofDetector{
		BaseDetector: analyzer.NewBaseDetector(5 * time.Minute),
		ipToMacs:     make(map[string]map[string]bool),
	}
}

func (d *ARPSpoofDetector) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" || packet.SrcMAC == "" {
		return alerts
	}

	macLower := strings.ToLower(packet.SrcMAC)
	if macLower == "ff:ff:ff:ff:ff:ff" || macLower == "00:00:00:00:00:00" {
		return alerts
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.ipToMacs[packet.SrcIP]; !exists {
		d.ipToMacs[packet.SrcIP] = make(map[string]bool)
	}
	d.ipToMacs[packet.SrcIP][macLower] = true

	if len(d.ipToMacs[packet.SrcIP]) > 1 {
		if d.ShouldAlert("ARP_SPOOF", packet.SrcIP) {
			macs := make([]string, 0, len(d.ipToMacs[packet.SrcIP]))
			for m := range d.ipToMacs[packet.SrcIP] {
				macs = append(macs, m)
			}
			alerts = append(alerts, analyzer.Alert{
				Timestamp: packet.Timestamp.Format(time.RFC3339),
				Level:     analyzer.LevelCritical,
				Type:      "ARP_SPOOF",
				Message:   "Address conflict: IP " + packet.SrcIP + " seen on different devices: " + strings.Join(macs, ", "),
				SrcIP:     packet.SrcIP,
			})
		}
	}

	return alerts
}
