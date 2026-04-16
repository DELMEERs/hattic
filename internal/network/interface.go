package network

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

type Interface struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IPAddresses []string `json:"ip_addresses"`
}

func GetInterfaces() ([]Interface, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find interfaces: %w", err)
	}

	var interfaces []Interface
	for _, device := range devices {
		var ips []string
		for _, addr := range device.Addresses {
			ips = append(ips, addr.IP.String())
		}
		interfaces = append(interfaces, Interface{
			Name:        device.Name,
			Description: device.Description,
			IPAddresses: ips,
		})
	}
	return interfaces, nil
}
