package network

import (
	"runtime"
	"strings"
	"time"

	"github.com/google/gopacket/pcap"
)

type SystemStatus struct {
	Status       string `json:"status"`
	Platform     string `json:"platform"`
	MissingDep   string `json:"missingDep"`
	Instructions string `json:"instructions"`
}

func GetSystemStatus() SystemStatus {
	platform := runtime.GOOS
	status := SystemStatus{
		Status:   "OK",
		Platform: platform,
	}

	v := pcap.Version()
	if v == "" || strings.Contains(strings.ToLower(v), "not found") {
		status.Status = "ERROR"
		if platform == "windows" {
			status.MissingDep = "Npcap"
			status.Instructions = "Npcap is required for packet capturing on Windows."
		} else {
			status.MissingDep = "libpcap"
			status.Instructions = "libpcap is required for packet capturing on Linux. Install it via your package manager."
		}
		return status
	}

	if platform == "linux" {
		handle, err := pcap.OpenLive("lo", 1024, false, 10*time.Millisecond)
		if err != nil {
			status.Status = "ERROR"
			status.MissingDep = "CAP_NET_RAW"
			status.Instructions = "sudo setcap cap_net_raw,cap_net_admin=eip ./hattic"
			return status
		}
		handle.Close()
	}

	return status
}