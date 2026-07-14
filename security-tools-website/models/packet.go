package models

type PacketInfo struct {
	Protocol string `json:"protocol"`

	SrcIP        string `json:"srcIp,omitempty"`
	DstIP        string `json:"dstIp,omitempty"`
	SrcPort      uint16 `json:"srcPort,omitempty"`
	DstPort      uint16 `json:"dstPort,omitempty"`
	Flags        string `json:"flags,omitempty"`
	Length       int    `json:"length"`
	HexDump      string `json:"hexDump,omitempty"`
	TypeCode     uint16 `json:"typeCode,omitempty"`
	SrcMAC       string `json:"srcMac,omitempty"`
	DstMAC       string `json:"dstMac,omitempty"`
	ARPOperation string `json:"arpOperation,omitempty"`
	ARPSenderIP  string `json:"arpSenderIp,omitempty"`
	ARPTargetIP  string `json:"arpTargetIp,omitempty"`
}

type StartCaptureRequest struct {
	Device string `json:"device"`
}

type DeviceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
