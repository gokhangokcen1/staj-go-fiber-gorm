package capture

import (
	"context"
	"fmt"
	"net"

	"github.com/gokhangokcen1/subnet-backend/models"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func StartCapture(deviceName string) (chan models.PacketInfo, context.CancelFunc, error) {
	handle, err := pcap.OpenLive(deviceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, nil, err
	}

	if err := handle.SetBPFFilter("tcp or udp or arp or icmp"); err != nil {
		handle.Close()
		return nil, nil, err
	}

	packetChan := make(chan models.PacketInfo, 200)
	ctx, cancel := context.WithCancel(context.Background())
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		defer handle.Close()
		defer close(packetChan)

		for {
			select {
			case <-ctx.Done():
				return
			case packet, ok := <-packetSource.Packets():
				if !ok {
					return
				}
				info, matched := parsePacket(packet)
				if !matched {
					continue
				}
				// info := parsePacket(packet)
				select {
				case packetChan <- info:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return packetChan, cancel, nil
}

func ListDevices() ([]pcap.Interface, error) {
	return pcap.FindAllDevs()
}

func parsePacket(packet gopacket.Packet) (models.PacketInfo, bool) {
	info := models.PacketInfo{}

	if arpLayer, ok := packet.Layer(layers.LayerTypeARP).(*layers.ARP); ok {
		info.Protocol = "ARP"
		info.SrcMAC = net.HardwareAddr(arpLayer.SourceHwAddress).String()
		info.DstMAC = net.HardwareAddr(arpLayer.DstHwAddress).String()
		info.ARPSenderIP = net.IP(arpLayer.SourceProtAddress).String()
		info.ARPTargetIP = net.IP(arpLayer.DstProtAddress).String()

		if arpLayer.Operation == layers.ARPRequest {
			info.ARPOperation = "Request"
		} else {
			info.ARPOperation = "Reply"
		}
		return info, true
	}

	ipLayer, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	if !ok {
		return info, false
	}
	info.SrcIP = ipLayer.SrcIP.String()
	info.DstIP = ipLayer.DstIP.String()

	if tcpLayer, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
		info.Protocol = "TCP"
		info.SrcPort = uint16(tcpLayer.SrcPort)
		info.DstPort = uint16(tcpLayer.DstPort)
		info.Flags = buildTCPFlags(tcpLayer)
	} else if udpLayer, ok := packet.Layer(layers.LayerTypeUDP).(*layers.UDP); ok {
		info.Protocol = "UDP"
		info.SrcPort = uint16(udpLayer.SrcPort)
		info.DstPort = uint16(udpLayer.DstPort)
	} else if icmpLayer, ok := packet.Layer(layers.LayerTypeICMPv4).(*layers.ICMPv4); ok {
		info.Protocol = "ICMP"
		info.TypeCode = uint16(icmpLayer.TypeCode)
	} else {
		return info, false
	}

	if appLayer := packet.ApplicationLayer(); appLayer != nil {
		payload := appLayer.Payload()
		info.Length = len(payload)
		info.HexDump = hexDump(payload)
	}

	return info, true
}

func buildTCPFlags(tcp *layers.TCP) string {
	flags := ""
	if tcp.SYN {
		flags += "SYN "
	}
	if tcp.ACK {
		flags += "ACK "
	}
	if tcp.FIN {
		flags += "FIN "
	}
	if tcp.RST {
		flags += "RST "
	}
	return flags
}

func hexDump(data []byte) string {
	var result string
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		hexPart := ""
		asciiPart := ""
		for _, b := range chunk {
			hexPart += fmt.Sprintf("%02x ", b)
			if b >= 32 && b <= 126 {
				asciiPart += string(b)
			} else {
				asciiPart += "."
			}
		}
		result += fmt.Sprintf("%-48s %s\n", hexPart, asciiPart)
	}
	return result
}
