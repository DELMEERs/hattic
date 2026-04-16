package config

import (
	"encoding/json"
	"os"
	"sync"
)

type AppConfig struct {
	InterfaceName string `json:"interface_name"`
	Promiscuous   bool   `json:"promiscuous"`
	SnapLen       int32  `json:"snap_len"`
}

type Manager struct {
	config AppConfig
	path   string
	mu     sync.RWMutex
}

func NewManager(path string) *Manager {
	return &Manager{
		path: path,
		config: AppConfig{
			Promiscuous: true,
			SnapLen:     65536,
		},
	}
}

func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return m.save()
		}
		return err
	}
	return json.Unmarshal(data, &m.config)
}

func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.save()
}

func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0644)
}

func (m *Manager) GetConfig() AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

func (m *Manager) SetConfig(config AppConfig) error {
	m.mu.Lock()
	m.config = config
	m.mu.Unlock()
	return m.Save()
}
