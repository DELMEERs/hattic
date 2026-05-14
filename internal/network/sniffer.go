package network

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"hattic/internal/analyzer"
	"hattic/internal/config"
)

type Sniffer struct {
	configManager *config.Manager
	analyzerMgr   *analyzer.Manager
	cancelFunc    context.CancelFunc
}

func NewSniffer(cfg *config.Manager, analyzerMgr *analyzer.Manager) *Sniffer {
	return &Sniffer{
		configManager: cfg,
		analyzerMgr:   analyzerMgr,
	}
}

func (s *Sniffer) Start() error {
	cfg := s.configManager.GetConfig()
	iface := cfg.InterfaceName

	// If no interface is specified, attempt to auto-detect the first non-loopback interface
	// that has at least one IP address assigned.
	if iface == "" {
		devices, err := pcap.FindAllDevs()
		if err != nil {
			return fmt.Errorf("failed to find devices for auto-detect: %w", err)
		}
		for _, dev := range devices {
			isLoopback := false
			for _, addr := range dev.Addresses {
				if addr.IP.IsLoopback() {
					isLoopback = true
					break
				}
			}
			if !isLoopback && len(dev.Addresses) > 0 {
				iface = dev.Name
				break
			}
		}
	}

	if iface == "" {
		return fmt.Errorf("no suitable network interface found. Please configure one manually.")
	}

	// OpenLive opens the device for real-time packet capture.
	// We use pcap.BlockForever to ensure the handle stays open until context is cancelled.
	handle, err := pcap.OpenLive(iface, cfg.SnapLen, cfg.Promiscuous, pcap.BlockForever)
	if err != nil {
		log.Printf("ERROR: Failed to open interface %s: %v", iface, err)
		// Check for common permission errors and provide actionable feedback.
		// On Linux, the user likely needs CAP_NET_RAW.
		if strings.Contains(err.Error(), "Permission denied") || strings.Contains(err.Error(), "socket: operation not permitted") {
			return fmt.Errorf("PERMISSION_DENIED: root or CAP_NET_RAW required to sniff on %s", iface)
		}
		return fmt.Errorf("failed to open device %s: %w", iface, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		defer handle.Close()
		log.Printf("Started sniffing on interface %s", iface)

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping sniffer...")
				return
			case packet, ok := <-packetSource.Packets():
				if !ok {
					return
				}
				s.processPacket(packet)
			}
		}
	}()

	return nil
}

func (s *Sniffer) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

var (
	tcpProtocols = map[uint16]string{
		80:  "HTTP",
		443: "HTTPS",
		22:  "SSH",
	}
	udpProtocols = map[uint16]string{
		53:   "DNS",
		5353: "MDNS",
	}
)

func (s *Sniffer) processPacket(packet gopacket.Packet) {
	info := &analyzer.PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Length:    packet.Metadata().Length,
	}

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth := ethLayer.(*layers.Ethernet)
		info.SrcMAC = eth.SrcMAC.String()
		info.DstMAC = eth.DstMAC.String()
	}

	if netLayer := packet.NetworkLayer(); netLayer != nil {
		info.SrcIP = netLayer.NetworkFlow().Src().String()
		info.DstIP = netLayer.NetworkFlow().Dst().String()

		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			ipv4 := ipv4Layer.(*layers.IPv4)
			info.TTL = ipv4.TTL
		}
	}

	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		info.Protocol = transportLayer.LayerType().String()

		switch t := transportLayer.(type) {
		case *layers.TCP:
			info.SrcPort = uint16(t.SrcPort)
			info.DstPort = uint16(t.DstPort)
			if proto, ok := tcpProtocols[info.DstPort]; ok {
				info.Protocol = proto
			} else if proto, ok := tcpProtocols[info.SrcPort]; ok {
				info.Protocol = proto
			}
		case *layers.UDP:
			info.SrcPort = uint16(t.SrcPort)
			info.DstPort = uint16(t.DstPort)
			if proto, ok := udpProtocols[info.DstPort]; ok {
				info.Protocol = proto
			} else if proto, ok := udpProtocols[info.SrcPort]; ok {
				info.Protocol = proto
			}
		}
	}

	if info.Protocol == "MDNS" {
		if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
			dns := dnsLayer.(*layers.DNS)
			for _, q := range dns.Questions {
				info.Hostname = string(q.Name)
				break
			}
		}
	}

	if info.Protocol == "" {
		if netLayer := packet.NetworkLayer(); netLayer != nil {
			info.Protocol = netLayer.LayerType().String()
		}
	}

	info.Protocol = strings.ToUpper(info.Protocol)
	s.analyzerMgr.ProcessPacket(info)
}
