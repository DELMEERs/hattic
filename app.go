package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"hattic/internal/analyzer"
	"hattic/internal/analyzer/detectors"
	"hattic/internal/config"
	"hattic/internal/network"
)

type App struct {
	ctx           context.Context
	configMgr     *config.Manager
	analyzerMgr   *analyzer.Manager
	sniffer       *network.Sniffer
	alertChan     chan analyzer.Alert
	stats         *Stats
	mu            sync.Mutex
	isSniffing    bool
}

type Stats struct {
	TotalPackets uint64            `json:"total_packets"`
	TotalAlerts  uint64            `json:"total_alerts"`
	ProtocolDist map[string]uint64 `json:"protocol_dist"`
}

func NewApp() *App {
	alertChan := make(chan analyzer.Alert, 100)
	configMgr := config.NewManager("config.json")
	_ = configMgr.Load()

	analyzerMgr := analyzer.NewManager(alertChan)
	analyzerMgr.RegisterDetector(detectors.NewARPSpoofDetector())
	analyzerMgr.RegisterDetector(detectors.NewMDNSScanner())
	analyzerMgr.RegisterDetector(detectors.NewPortScannerDetector(50))
	analyzerMgr.RegisterDetector(detectors.NewSuspiciousProtocolDetector())
	analyzerMgr.RegisterDetector(detectors.NewTrafficFloodDetector(1000, 5000, 10000))
	analyzerMgr.RegisterDetector(detectors.NewUnusualTTLDetector())

	app := &App{
		configMgr:   configMgr,
		analyzerMgr: analyzerMgr,
		alertChan:   alertChan,
		stats: &Stats{
			ProtocolDist: make(map[string]uint64),
		},
	}

	analyzerMgr.OnPacket = func(p *analyzer.PacketInfo) {
		app.mu.Lock()
		app.stats.TotalPackets++
		app.stats.ProtocolDist[p.Protocol]++
		app.mu.Unlock()

		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "packet", p)
		}
	}

	app.sniffer = network.NewSniffer(configMgr, analyzerMgr)

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.listenForAlerts()
	go func() {
		time.Sleep(500 * time.Millisecond)
		_ = a.StartSniffer()
	}()
}

func (a *App) listenForAlerts() {
	for alert := range a.alertChan {
		a.mu.Lock()
		a.stats.TotalAlerts++
		a.mu.Unlock()
		runtime.EventsEmit(a.ctx, "alert", alert)
	}
}


func (a *App) GetInterfaces() ([]network.NetworkInterface, error) {
	return network.GetInterfaces()
}

func (a *App) GetConfig() config.AppConfig {
	return a.configMgr.GetConfig()
}

func (a *App) SaveConfig(cfg config.AppConfig) error {
	return a.configMgr.SetConfig(cfg)
}

type HealthStatus struct {
	IsRoot      bool   `json:"is_root"`
	PcapVersion string `json:"pcap_version"`
	CanSniff    bool   `json:"can_sniff"`
	Error       string `json:"error"`
}

func (a *App) HealthCheck() HealthStatus {
	status := HealthStatus{
		IsRoot:      os.Geteuid() == 0,
		PcapVersion: pcap.Version(),
		CanSniff:    true,
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		status.CanSniff = false
		status.Error = "Failed to list devices: " + err.Error()
		return status
	}

	if len(devices) == 0 {
		status.CanSniff = false
		status.Error = "No network devices found"
		return status
	}

	return status
}

func (a *App) StartSniffer() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isSniffing {
		return fmt.Errorf("sniffer is already running")
	}

	err := a.sniffer.Start()
	if err != nil {
		runtime.EventsEmit(a.ctx, "backend-error", err.Error())
		return err
	}

	a.isSniffing = true
	runtime.EventsEmit(a.ctx, "sniffer_status", true)
	return nil
}

func (a *App) StopSniffer() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isSniffing {
		return
	}

	a.sniffer.Stop()
	a.isSniffing = false
	runtime.EventsEmit(a.ctx, "sniffer_status", false)
}

func (a *App) GetStats() *Stats {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stats
}

func (a *App) GetIsSniffing() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isSniffing
}

func (a *App) TriggerTestAlert() {
	alert := analyzer.Alert{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     analyzer.LevelCritical,
		Type:      "TEST_ANOMALY",
		Message:   "This is a simulated critical threat detected by the system.",
		SrcIP:     "1.2.3.4",
	}
	a.mu.Lock()
	a.stats.TotalAlerts++
	a.mu.Unlock()
	runtime.EventsEmit(a.ctx, "alert", alert)
}