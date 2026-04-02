package sniffer

import (
	"fmt"
	"net"
	"time"

	"hattic/internal/db"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type FlowKey struct {
	SrcIP    string
	DstIP    string
	SrcMAC   string
	DstMAC   string
	Protocol string
}

type Sniffer struct {
	Device      string
	SnapshotLen int32
	Promiscuous bool
	Timeout     time.Duration
	LogChan     chan<- db.TrafficLog
}

func NewSniffer(device string, logChan chan<- db.TrafficLog) *Sniffer {
	return &Sniffer{
		Device:      device,
		SnapshotLen: 1024,
		Promiscuous: false,
		Timeout:     pcap.BlockForever,
		LogChan:     logChan,
	}
}

func (s *Sniffer) Start() error {
	handle, err := pcap.OpenLive(s.Device, s.SnapshotLen, s.Promiscuous, s.Timeout)
	if err != nil {
		return fmt.Errorf("error opening device %s: %v", s.Device, err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	flows := make(map[FlowKey]*db.TrafficLog)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case packet := <-packetSource.Packets():
			if packet == nil {
				return nil
			}
			s.processPacket(packet, flows)
		case <-ticker.C:
			s.flushFlows(flows)
		}
	}
}

func (s *Sniffer) processPacket(packet gopacket.Packet, flows map[FlowKey]*db.TrafficLog) {
	var srcIP, dstIP, srcMAC, dstMAC, protocol string
	var payloadSize int

	// Ethernet Layer
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer != nil {
		eth := ethLayer.(*layers.Ethernet)
		srcMAC = eth.SrcMAC.String()
		dstMAC = eth.DstMAC.String()
	}

	// IPv4 Layer
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip := ipLayer.(*layers.IPv4)
		srcIP = ip.SrcIP.String()
		dstIP = ip.DstIP.String()
		protocol = ip.Protocol.String()
		payloadSize = len(ip.Payload)
	}

	// ARP Layer
	arpLayer := packet.Layer(layers.LayerTypeARP)
	if arpLayer != nil {
		arp := arpLayer.(*layers.ARP)
		srcIP = net.IP(arp.SourceProtAddress).String()
		dstIP = net.IP(arp.DstProtAddress).String()
		srcMAC = net.HardwareAddr(arp.SourceHwAddress).String()
		dstMAC = net.HardwareAddr(arp.DstHwAddress).String()
		protocol = "ARP"
		payloadSize = len(arp.Payload)
	}

	// UDP Layer (mDNS)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		if udp.DstPort == 5353 || udp.SrcPort == 5353 {
			protocol = "mDNS"
		}
	}

	if srcIP == "" && protocol == "" {
		return
	}

	key := FlowKey{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		SrcMAC:   srcMAC,
		DstMAC:   dstMAC,
		Protocol: protocol,
	}

	if log, ok := flows[key]; ok {
		log.PacketCount++
		log.PayloadSize += payloadSize
		log.Timestamp = time.Now()
	} else {
		flows[key] = &db.TrafficLog{
			Timestamp:   time.Now(),
			SrcIP:       srcIP,
			DstIP:       dstIP,
			SrcMAC:      srcMAC,
			DstMAC:      dstMAC,
			Protocol:    protocol,
			PayloadSize: payloadSize,
			PacketCount: 1,
		}
	}
}

func (s *Sniffer) flushFlows(flows map[FlowKey]*db.TrafficLog) {
	for key, logEntry := range flows {
		s.LogChan <- *logEntry
		delete(flows, key)
	}
}
