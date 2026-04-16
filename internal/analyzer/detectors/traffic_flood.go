package detectors

import (
	"fmt"
	"hattic/internal/analyzer"
	"sync"
	"time"
)

type TrafficFloodDetector struct {
	*analyzer.BaseDetector
	infoThreshold     int
	warningThreshold  int
	criticalThreshold int
	ipPacketCount     map[string]int
	mu                sync.Mutex
	lastReset         time.Time
}

func NewTrafficFloodDetector(info, warning, critical int) *TrafficFloodDetector {
	return &TrafficFloodDetector{
		BaseDetector:      analyzer.NewBaseDetector(1 * time.Minute),
		infoThreshold:     info,
		warningThreshold:  warning,
		criticalThreshold: critical,
		ipPacketCount:     make(map[string]int),
		lastReset:         time.Now(),
	}
}

func (d *TrafficFloodDetector) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" {
		return alerts
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// reset counters every minute idk
	if time.Since(d.lastReset) > time.Minute {
		d.ipPacketCount = make(map[string]int)
		d.lastReset = time.Now()
	}

	d.ipPacketCount[packet.SrcIP]++
	pCount := d.ipPacketCount[packet.SrcIP]

	var level analyzer.Level
	desc := ""
	matched := false

	if pCount >= d.criticalThreshold {
		level = analyzer.LevelCritical
		desc = "Potential DoS attack detected!"
		matched = true
	} else if pCount >= d.warningThreshold {
		level = analyzer.LevelWarning
		desc = "Extreme network activity (possible flood)."
		matched = true
	} else if pCount >= d.infoThreshold {
		level = analyzer.LevelInfo
		desc = "High network activity (streaming/downloading)."
		matched = true
	}

	if matched {
		if d.ShouldAlert("TRAFFIC_FLOOD", packet.SrcIP) {
			alerts = append(alerts, analyzer.Alert{
				Timestamp: packet.Timestamp.Format(time.RFC3339),
				Level:     level,
				Type:      "TRAFFIC_FLOOD",
				Message:   fmt.Sprintf("[%s] %s (%d packets)", packet.SrcIP, desc, pCount),
				SrcIP:     packet.SrcIP,
			})
		}
	}

	return alerts
}
