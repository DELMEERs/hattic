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

	if cfg.InterfaceName == "" {
		return fmt.Errorf("interface name is not configured")
	}

	handle, err := pcap.OpenLive(cfg.InterfaceName, cfg.SnapLen, cfg.Promiscuous, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", cfg.InterfaceName, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		defer handle.Close()
		log.Printf("Started sniffing on interface %s", cfg.InterfaceName)

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

func (s *Sniffer) processPacket(packet gopacket.Packet) {
	info := &analyzer.PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Length:    packet.Metadata().Length,
	}

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth, _ := ethLayer.(*layers.Ethernet)
		info.SrcMAC = eth.SrcMAC.String()
		info.DstMAC = eth.DstMAC.String()
	}

	if netLayer := packet.NetworkLayer(); netLayer != nil {
		info.SrcIP = netLayer.NetworkFlow().Src().String()
		info.DstIP = netLayer.NetworkFlow().Dst().String()

		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			ipv4, _ := ipv4Layer.(*layers.IPv4)
			info.TTL = ipv4.TTL
		}
	}

	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		info.Protocol = transportLayer.LayerType().String()

		switch t := transportLayer.(type) {
		case *layers.TCP:
			info.SrcPort = uint16(t.SrcPort)
			info.DstPort = uint16(t.DstPort)
			if info.DstPort == 80 || info.SrcPort == 80 {
				info.Protocol = "HTTP"
			} else if info.DstPort == 443 || info.SrcPort == 443 {
				info.Protocol = "HTTPS"
			} else if info.DstPort == 22 || info.SrcPort == 22 {
				info.Protocol = "SSH"
			}
		case *layers.UDP:
			info.SrcPort = uint16(t.SrcPort)
			info.DstPort = uint16(t.DstPort)
			if info.DstPort == 53 || info.SrcPort == 53 {
				info.Protocol = "DNS"
			} else if info.DstPort == 5353 || info.SrcPort == 5353 {
				info.Protocol = "MDNS"
			}
		}
	}

	if info.Protocol == "MDNS" {
		if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
			dns, _ := dnsLayer.(*layers.DNS)
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
