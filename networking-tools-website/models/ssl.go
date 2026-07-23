package models

type SslCheckRequest struct {
	Website string `json:"website"`
	Port    int    `json:"port"`
}

type GeneralInformation struct {
	ResolvesTo       string
	ExpirationDate   string
	VendorSigned     bool
	HostnameMatches  bool
	KeyLength        int
	ServerType       string
	RevocationStatus string
}

type IssuedFor struct {
	CommonName       string
	SAN              []string
	Organization     string
	OrganizationUnit string
	Country          string
	State            string
	Locality         string
	Address          string
}

type IssuedBy struct {
	CommonName       string
	Organization     string
	OrganizationUnit string
	Country          string
	State            string
	Locality         string
}

type ChainCert struct {
	Issuer             string
	CommonName         string
	Organization       string
	Issued             string
	Expires            string
	SerialNumber       string
	SignatureAlgorithm string
	FingerprintSHA1    string
	FingerprintMD5     string
	PEM                string
}

type ChainDetails struct {
	Certs []ChainCert
}

type SSLReport struct {
	General GeneralInformation
	For     IssuedFor
	By      IssuedBy
	Chain   ChainDetails
}
