package models

// SubnetRequest, kullanicidan (Vue formundan) gelen istegi temsil eder.
type SubnetRequest struct {
	IP   string `json:"ip"`
	CIDR int    `json:"cidr"`
}

// SubnetResponse, hesaplama sonucunu istemciye donen yapidir.
type SubnetResponse struct {
	Adres              string  `json:"adres"`
	NetworkMask        string  `json:"networkMask"`
	WildcardMask       string  `json:"wildcardMask"`
	NetworkAddress     string  `json:"networkAddress"`
	BroadcastAddress   string  `json:"broadcastAddress"`
	Hostmin            string  `json:"hostmin"`
	Hostmax            string  `json:"hostmax"`
	HostsPerNet        float64 `json:"hostsPerNet"`
	AddressBinary      string  `json:"addressBinary"`
	MaskBinary         string  `json:"maskBinary"`
	NetworkBinary      string  `json:"networkBinary"`
	BroadcastBinary    string  `json:"broadcastBinary"`
	HostMinBinary      string  `json:"hostMinBinary"`
	HostMaxBinary      string  `json:"hostMaxBinary"`
	WildcardMaskBinary string  `json:"wildcardMaskBinary"`
	Class              string  `json:"class"`
}
