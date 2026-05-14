package analyzer

import (
	"sync"
)

type Manager struct {
	detectors []Detector
	mu        sync.RWMutex
	alertChan chan<- Alert
	OnPacket  func(*PacketInfo)
}

func NewManager(alertChan chan<- Alert) *Manager {
	return &Manager{
		detectors: make([]Detector, 0),
		alertChan: alertChan,
	}
}

func (m *Manager) RegisterDetector(d Detector) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.detectors = append(m.detectors, d)
}

func (m *Manager) ProcessPacket(packet *PacketInfo) {
	if m.OnPacket != nil {
		m.OnPacket(packet)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.detectors {
		alerts := d.Analyze(packet)
		for _, alert := range alerts {
			select {
			case m.alertChan <- alert:
			default:
				// Discard alert if channel is full to avoid blocking the capture loop.
				// In a production environment, consider buffering or logging dropped alerts.
			}
		}
	}
}
