package models

type PortCheckRequest struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type PortCheckResponse struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Acik  bool   `json:"acik"`
	Detay string `json:"detay"`
}
