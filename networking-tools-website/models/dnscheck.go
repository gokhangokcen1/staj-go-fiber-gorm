package models

type DNSCheckRequest struct {
	Domain     string `json:"domain"`
	RecordType string `json:"recordType"`
}

type DNSPropagationResult struct {
	Server    string   `json:"server"`
	Location  string   `json:"location"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Status    string   `json:"status"`
	Answers   []string `json:"answers,omitempty"`
	Error     string   `json:"error,omitempty"`
}

type DNSPropagationResponse struct {
	Domain     string                 `json:"domain"`
	RecordType string                 `json:"recordType"`
	Results    []DNSPropagationResult `json:"results"`
}
