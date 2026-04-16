package detectors

import (
	"fmt"
	"hattic/internal/analyzer"
	"strings"
	"time"
)

type SuspiciousProtocolDetector struct {
	*analyzer.BaseDetector
	commonPorts        map[uint16]map[string]bool
	suspiciousMappings map[uint16]map[string]bool
}

func NewSuspiciousProtocolDetector() *SuspiciousProtocolDetector {
	d := &SuspiciousProtocolDetector{
		BaseDetector:       analyzer.NewBaseDetector(5 * time.Minute),
		commonPorts:        make(map[uint16]map[string]bool),
		suspiciousMappings: make(map[uint16]map[string]bool),
	}

	d.commonPorts[80] = map[string]bool{"TCP": true, "HTTP": true}
	d.commonPorts[443] = map[string]bool{"TCP": true, "HTTPS": true}
	d.commonPorts[53] = map[string]bool{"UDP": true, "TCP": true, "DNS": true}
	d.commonPorts[22] = map[string]bool{"TCP": true, "SSH": true}
	d.commonPorts[25] = map[string]bool{"TCP": true, "SMTP": true}
	d.commonPorts[5353] = map[string]bool{"UDP": true, "MDNS": true}

	d.suspiciousMappings[80] = map[string]bool{"SSH": true}
	d.suspiciousMappings[22] = map[string]bool{"UDP": true}

	return d
}

func (d *SuspiciousProtocolDetector) Analyze(packet *analyzer.PacketInfo) []analyzer.Alert {
	var alerts []analyzer.Alert

	if packet.SrcIP == "" || packet.DstPort == 0 || packet.Protocol == "" {
		return alerts
	}

	proto := strings.ToUpper(packet.Protocol)
	isSuspicious := false
	reason := ""

	if protos, ok := d.suspiciousMappings[packet.DstPort]; ok && protos[proto] {
		isSuspicious = true
		reason = fmt.Sprintf("Unusual protocol %s on port %d", proto, packet.DstPort)
	}

	if isSuspicious {
		if d.ShouldAlert("SUSPICIOUS_PROTOCOL", packet.SrcIP) {
			alerts = append(alerts, analyzer.Alert{
				Timestamp: packet.Timestamp.Format(time.RFC3339),
				Level:     analyzer.LevelWarning,
				Type:      "SUSPICIOUS_PROTOCOL",
				Message:   fmt.Sprintf("[%s] %s", packet.SrcIP, reason),
				SrcIP:     packet.SrcIP,
			})
		}
	}

	return alerts
}
