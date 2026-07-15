package models

type ScanRequest struct {
	IP   string `json:"ip"`
	CIDR int    `json:"cidr"`
}
