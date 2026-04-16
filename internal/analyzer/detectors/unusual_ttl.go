package detectors

import (
	"fmt"
	"hattic/internal/analyzer"
	"strings"
	"time"
)

type UnusualTTLDetector struct {
	*analyzer.BaseDetector
	standardTTLs map[uint8]bool
}

func NewUnusualTTLDetector() *UnusualTTLDetector {
	return &UnusualTTLDetector{
		BaseDetector: analyzer.NewBaseDetector(5 * time.Minute),
		standardTTLs: map[uint8]bool{
			64:  true,
			128: true,
			255: true,
		},
	}
}

func (d *UnusualTTLDetector) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" || packet.TTL == 0 {
		return alerts
	}

	if d.standardTTLs[packet.TTL] {
		return alerts
	}

	isLocal := strings.HasPrefix(packet.SrcIP, "192.168.") ||
		strings.HasPrefix(packet.SrcIP, "10.") ||
		strings.HasPrefix(packet.SrcIP, "172.16.")

	if isLocal {
		return alerts
	}

	if packet.TTL >= 10 && packet.TTL <= 30 {
		if d.ShouldAlert("UNUSUAL_TTL", packet.SrcIP) {
			alerts = append(alerts, analyzer.Alert{
				Timestamp: packet.Timestamp.Format(time.RFC3339),
				Level:     analyzer.LevelInfo,
				Type:      "UNUSUAL_TTL",
				Message:   fmt.Sprintf("External IP %s with unusual TTL: %d", packet.SrcIP, packet.TTL),
				SrcIP:     packet.SrcIP,
			})
		}
	}

	return alerts
}
