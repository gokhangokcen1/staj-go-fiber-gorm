package models

type WhoisRequest struct {
	Domain string `json:"domain"`
}

type WhoisContact struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Fax          string `json:"fax"`
	Email        string `json:"email"`
}

type WhoisRegistrar struct {
	Name       string `json:"name"`
	IANAID     string `json:"ianaId"`
	Email      string `json:"email"`
	AbuseEmail string `json:"abuseEmail"`
	AbusePhone string `json:"abusePhone"`
}

type WhoisResponse struct {
	Domain       string         `json:"domain"`
	RegisteredOn string         `json:"registeredOn"`
	ExpiresOn    string         `json:"expiresOn"`
	UpdatedOn    string         `json:"updatedOn"`
	Status       []string       `json:"status"`
	NameServers  []string       `json:"nameServers"`
	Registrar    WhoisRegistrar `json:"registrar"`
	Registrant   WhoisContact   `json:"registrant"`
	Technical    WhoisContact   `json:"technical"`
}
