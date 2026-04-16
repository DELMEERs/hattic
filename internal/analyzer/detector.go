package analyzer

import (
	"sync"
	"time"
)

type Detector interface {
	Analyze(packet *PacketInfo) []Alert
}

type BaseDetector struct {
	lastAlertTimes map[string]time.Time
	cooldown       time.Duration
	mu             sync.Mutex
}

func NewBaseDetector(cooldown time.Duration) *BaseDetector {
	return &BaseDetector{
		lastAlertTimes: make(map[string]time.Time),
		cooldown:       cooldown,
	}
}

func (b *BaseDetector) ShouldAlert(alertType string, srcIP string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	key := alertType + "|" + srcIP
	now := time.Now()

	if lastAlert, exists := b.lastAlertTimes[key]; exists {
		if now.Sub(lastAlert) < b.cooldown {
			return false
		}
	}

	b.lastAlertTimes[key] = now
	return true
}
