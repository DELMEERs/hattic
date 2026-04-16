package analyzer

import (
	"sync"
)

type Manager struct {
	detectors []Detector
	mu        sync.RWMutex
	alertChan chan<- Alert
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
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.detectors {
		alerts := d.Analyze(packet)
		for _, alert := range alerts {
			select {
			case m.alertChan <- alert:
			default:
				// if the channel is full log it or discard it. for now we are using an unbuffered channel and not blocking
				// in a real scenario I use a buffered channel or a context timeout
				// using select with the default value we discard the alert if the receiver is not ready
				// we can simply rely on a normal send through the channel if we prefer blocking but discarding is safer to prevent blocking the capture loop
				// i will assume the alertChan buffer is a reasonable size
			}
		}
	}
}
