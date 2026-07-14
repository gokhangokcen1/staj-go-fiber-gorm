package sslcheck

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors" // Sertifika gelmeme hatası için eklendi
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gokhangokcen1/subnet-backend/models"
	"golang.org/x/crypto/ocsp"
)

func normalizeWebsite(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")
	input = strings.TrimSuffix(input, "/")

	if idx := strings.Index(input, "/"); idx != -1 {
		input = input[:idx]
	}
	return input
}

func CheckSSL(rawWebsite string) (*models.SSLReport, error) {
	website := normalizeWebsite(rawWebsite)

	conn, err := tls.Dial("tcp", website+":443", &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         website,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connState := conn.ConnectionState()

	if len(connState.PeerCertificates) == 0 {
		return nil, errors.New("sunucu herhangi bir sertifika gondermedi")
	}

	leafCert := connState.PeerCertificates[0]

	var issuerCert *x509.Certificate
	if len(connState.PeerCertificates) > 1 {
		issuerCert = connState.PeerCertificates[1]
	}

	report := &models.SSLReport{}
	FillGeneralInformation(&report.General, leafCert, issuerCert, &connState, website)
	FillIssuedFor(&report.For, leafCert)
	FillIssuedBy(&report.By, leafCert)
	FillChainDetails(&report.Chain, &connState)

	return report, nil
}

func FillGeneralInformation(gi *models.GeneralInformation, cert *x509.Certificate, issuer *x509.Certificate, connState *tls.ConnectionState, website string) {

	gi.ResolvesTo = website

	gi.ExpirationDate = cert.NotAfter.Format(time.DateOnly)

	intermediates := x509.NewCertPool()
	for _, c := range connState.PeerCertificates[1:] {
		intermediates.AddCert(c)
	}

	roots, err := x509.SystemCertPool()
	if err != nil || roots == nil {
		roots = x509.NewCertPool()
	}

	opts := x509.VerifyOptions{
		DNSName:       "",
		Intermediates: intermediates,
		Roots:         roots,
	}

	_, errVerify := cert.Verify(opts)
	gi.VendorSigned = (errVerify == nil)

	errHost := cert.VerifyHostname(website)
	gi.HostnameMatches = (errHost == nil)

	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		gi.KeyLength = pub.N.BitLen()
	case *ecdsa.PublicKey:
		gi.KeyLength = pub.Params().BitSize
	case ed25519.PublicKey:
		gi.KeyLength = 256
	}

	gi.ServerType = "NA"
	gi.RevocationStatus = checkRevocation(cert, issuer)
}

func FillIssuedFor(f *models.IssuedFor, cert *x509.Certificate) {
	f.CommonName = cert.Subject.CommonName
	f.SAN = cert.DNSNames
	f.Organization = firstOrNA(cert.Subject.Organization)
	f.OrganizationUnit = firstOrNA(cert.Subject.OrganizationalUnit)
	f.Country = firstOrNA(cert.Subject.Country)
	f.State = firstOrNA(cert.Subject.Province)
	f.Locality = firstOrNA(cert.Subject.Locality)
	f.Address = firstOrNA(cert.Subject.StreetAddress)
}

func FillIssuedBy(b *models.IssuedBy, cert *x509.Certificate) {
	b.CommonName = cert.Issuer.CommonName
	b.Organization = firstOrNA(cert.Issuer.Organization)
	b.OrganizationUnit = firstOrNA(cert.Issuer.OrganizationalUnit)
	b.Country = firstOrNA(cert.Issuer.Country)
	b.State = firstOrNA(cert.Issuer.Province)
	b.Locality = firstOrNA(cert.Issuer.Locality)
}

func FillChainDetails(cd *models.ChainDetails, connState *tls.ConnectionState) {
	for _, cert := range connState.PeerCertificates {
		var oneCert models.ChainCert
		oneCert.Issuer = cert.Issuer.CommonName
		oneCert.CommonName = cert.Subject.CommonName
		oneCert.Organization = firstOrNA(cert.Subject.Organization)
		oneCert.Issued = cert.NotBefore.Format(time.DateOnly)
		oneCert.Expires = cert.NotAfter.Format(time.DateOnly)
		oneCert.SerialNumber = formatSerial(cert.SerialNumber)
		oneCert.SignatureAlgorithm = cert.SignatureAlgorithm.String()
		oneCert.FingerprintSHA1 = fmt.Sprintf("%x", sha1.Sum(cert.Raw))
		oneCert.FingerprintMD5 = fmt.Sprintf("%x", md5.Sum(cert.Raw))

		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		oneCert.PEM = string(pemBytes)

		cd.Certs = append(cd.Certs, oneCert)
	}
}

func firstOrNA(values []string) string {
	if len(values) > 0 && values[0] != "" {
		return values[0]
	}
	return "NA"
}

func formatSerial(sn *big.Int) string {
	bytes := sn.Bytes()
	parts := make([]string, len(bytes))
	for i, b := range bytes {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, ":")
}

func checkRevocation(leaf *x509.Certificate, issuer *x509.Certificate) string {
	if issuer == nil {
		return "Unknown"
	}

	if len(leaf.OCSPServer) == 0 {
		return "NA"
	}

	ocspServer := leaf.OCSPServer[0]

	ocspRequest, err := ocsp.CreateRequest(leaf, issuer, nil)
	if err != nil {
		return "Unknown"
	}

	httpResp, err := http.Post(ocspServer, "application/ocsp-request", bytes.NewReader(ocspRequest))
	if err != nil {
		return "Unknown"
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "Unknown"
	}

	ocspResponse, err := ocsp.ParseResponse(body, issuer)
	if err != nil {
		return "Unknown"
	}

	switch ocspResponse.Status {
	case ocsp.Good:
		return "Good"
	case ocsp.Revoked:
		return "Revoked"
	default:
		return "Unknown"
	}
}
