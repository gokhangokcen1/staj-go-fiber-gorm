package whois

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gokhangokcen1/subnet-backend/models"
)

const queryTimeout = 12 * time.Second

func Lookup(input string) (*models.WhoisResponse, error) {
	domain, err := normalizeDomain(input)
	if err != nil {
		return nil, err
	}
	server := registryServer(domain)
	if server == "" {
		tld := domain[strings.LastIndex(domain, ".")+1:]
		root, rootErr := query("whois.iana.org", tld)
		server = findValue(root, "refer")
		if rootErr != nil || server == "" {
			server = "whois." + tld
		}
	}
	registry, err := query(server, domain)
	if err != nil {
		return lookupRDAP(domain, err)
	}
	registrarRaw := ""
	if registrarServer := findValue(registry, "registrar whois server"); registrarServer != "" && !strings.EqualFold(registrarServer, server) {
		registrarRaw, _ = query(registrarServer, domain)
	}
	result := parse(domain, registry, registrarRaw)
	if result.RegisteredOn == "" && result.Registrar.Name == "" && len(result.NameServers) == 0 {
		return lookupRDAP(domain, fmt.Errorf("alan adi kaydi bulunamadi veya bu uzanti WHOIS sorgusunu desteklemiyor"))
	}
	return result, nil
}

func registryServer(domain string) string {
	switch {
	case strings.HasSuffix(domain, ".com"), strings.HasSuffix(domain, ".net"):
		return "whois.verisign-grs.com"
	case strings.HasSuffix(domain, ".org"):
		return "whois.publicinterestregistry.org"
	case strings.HasSuffix(domain, ".tr"):
		return "whois.trabis.gov.tr"
	case strings.HasSuffix(domain, ".edu"):
		return "whois.educause.net"
	case strings.HasSuffix(domain, ".gov"):
		return "whois.nic.gov"
	default:
		return ""
	}
}

func lookupRDAP(domain string, whoisErr error) (*models.WhoisResponse, error) {
	data, err := rdapData("https://rdap.org/domain/" + url.PathEscape(domain))
	if err != nil {
		return nil, fmt.Errorf("WHOIS sunucusuna ulasilamadi ve RDAP sorgusu da basarisiz oldu: %w", whoisErr)
	}
	result := &models.WhoisResponse{Domain: text(data["ldhName"])}
	if result.Domain == "" {
		result.Domain = domain
	}
	for _, event := range objects(data["events"]) {
		switch strings.ToLower(text(event["eventAction"])) {
		case "registration":
			result.RegisteredOn = formatDate(text(event["eventDate"]))
		case "expiration", "expiry":
			result.ExpiresOn = formatDate(text(event["eventDate"]))
		case "last changed", "last update":
			result.UpdatedOn = formatDate(text(event["eventDate"]))
		}
	}
	result.Status = stringsFrom(data["status"])
	for _, server := range objects(data["nameservers"]) {
		if name := text(server["ldhName"]); name != "" {
			result.NameServers = append(result.NameServers, strings.ToLower(name))
		}
	}
	applyEntities(result, objects(data["entities"]))
	for _, link := range objects(data["links"]) {
		if strings.EqualFold(text(link["rel"]), "related") {
			if details, err := rdapData(text(link["href"])); err == nil {
				applyEntities(result, objects(details["entities"]))
			}
			break
		}
	}
	sort.Strings(result.Status)
	sort.Strings(result.NameServers)
	return result, nil
}

func rdapData(endpoint string) (map[string]any, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/rdap+json, application/json")
	response, err := (&http.Client{Timeout: queryTimeout}).Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("RDAP sunucusu %s dondu", response.Status)
	}
	var data map[string]any
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func applyEntities(result *models.WhoisResponse, entities []map[string]any) {
	for _, entity := range entities {
		roles := stringsFrom(entity["roles"])
		if role(roles, "registrar") {
			result.Registrar = mergeRegistrar(result.Registrar, registrar(entity))
		}
		if role(roles, "registrant") {
			result.Registrant = mergeContact(result.Registrant, rdapContact(entity))
		}
		if role(roles, "technical") {
			result.Technical = mergeContact(result.Technical, rdapContact(entity))
		}
	}
}

func registrar(entity map[string]any) models.WhoisRegistrar {
	contact := rdapContact(entity)
	result := models.WhoisRegistrar{Name: contact.Name, Email: contact.Email}
	if contact.Organization != "" {
		result.Name = contact.Organization
	}
	for _, id := range objects(entity["publicIds"]) {
		if strings.Contains(strings.ToLower(text(id["type"])), "iana") {
			result.IANAID = text(id["identifier"])
		}
	}
	for _, child := range objects(entity["entities"]) {
		if role(stringsFrom(child["roles"]), "abuse") {
			abuse := rdapContact(child)
			result.AbuseEmail, result.AbusePhone = abuse.Email, abuse.Phone
		}
	}
	return result
}

func rdapContact(entity map[string]any) models.WhoisContact {
	result := models.WhoisContact{}
	vcard, _ := entity["vcardArray"].([]any)
	if len(vcard) < 2 {
		return result
	}
	properties, _ := vcard[1].([]any)
	for _, value := range properties {
		property, ok := value.([]any)
		if !ok || len(property) < 4 {
			continue
		}
		switch strings.ToLower(text(property[0])) {
		case "fn":
			result.Name = text(property[3])
		case "org":
			result.Organization = text(property[3])
		case "email":
			result.Email = text(property[3])
		case "contact-uri":
			if result.Email == "" {
				result.Email = text(property[3])
			}
		case "tel":
			if result.Phone == "" {
				result.Phone = text(property[3])
			} else if result.Fax == "" {
				result.Fax = text(property[3])
			}
		case "adr":
			parts, _ := property[3].([]any)
			if len(parts) >= 7 {
				result.Street, result.City, result.State, result.PostalCode, result.Country = text(parts[2]), text(parts[3]), text(parts[4]), text(parts[5]), text(parts[6])
			}
		}
	}
	return result
}

func mergeContact(old, newer models.WhoisContact) models.WhoisContact {
	if newer.Name != "" {
		old.Name = newer.Name
	}
	if newer.Organization != "" {
		old.Organization = newer.Organization
	}
	if newer.Street != "" {
		old.Street = newer.Street
	}
	if newer.City != "" {
		old.City = newer.City
	}
	if newer.State != "" {
		old.State = newer.State
	}
	if newer.PostalCode != "" {
		old.PostalCode = newer.PostalCode
	}
	if newer.Country != "" {
		old.Country = newer.Country
	}
	if newer.Phone != "" {
		old.Phone = newer.Phone
	}
	if newer.Fax != "" {
		old.Fax = newer.Fax
	}
	if newer.Email != "" {
		old.Email = newer.Email
	}
	return old
}

func mergeRegistrar(old, newer models.WhoisRegistrar) models.WhoisRegistrar {
	if newer.Name != "" {
		old.Name = newer.Name
	}
	if newer.IANAID != "" {
		old.IANAID = newer.IANAID
	}
	if newer.Email != "" {
		old.Email = newer.Email
	}
	if newer.AbuseEmail != "" {
		old.AbuseEmail = newer.AbuseEmail
	}
	if newer.AbusePhone != "" {
		old.AbusePhone = newer.AbusePhone
	}
	return old
}

func objects(value any) []map[string]any {
	values, _ := value.([]any)
	result := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if object, ok := value.(map[string]any); ok {
			result = append(result, object)
		}
	}
	return result
}
func stringsFrom(value any) []string {
	values, _ := value.([]any)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if stringValue := text(value); stringValue != "" {
			result = append(result, stringValue)
		}
	}
	return result
}
func text(value any) string { result, _ := value.(string); return result }
func role(roles []string, wanted string) bool {
	for _, value := range roles {
		if strings.EqualFold(value, wanted) {
			return true
		}
	}
	return false
}

func normalizeDomain(input string) (string, error) {
	raw := strings.TrimSpace(strings.ToLower(input))
	if raw == "" {
		return "", fmt.Errorf("alan adi zorunludur")
	}
	if !strings.Contains(raw, "://") {
		raw = "//" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("gecerli bir alan adi girin")
	}
	host := strings.TrimSuffix(u.Hostname(), ".")
	if host == "" || net.ParseIP(host) != nil || !strings.Contains(host, ".") {
		return "", fmt.Errorf("gecerli bir alan adi girin")
	}
	for _, label := range strings.Split(host, ".") {
		if label == "" || len(label) > 63 {
			return "", fmt.Errorf("gecerli bir alan adi girin")
		}
		for _, char := range label {
			if !(char == '-' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9') {
				return "", fmt.Errorf("gecerli bir alan adi girin")
			}
		}
	}
	return host, nil
}

func query(server, domain string) (string, error) {
	server = strings.TrimSpace(strings.TrimPrefix(server, "whois://"))
	if !strings.Contains(server, ":") {
		server += ":43"
	}
	conn, err := net.DialTimeout("tcp", server, queryTimeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(queryTimeout))
	if _, err := conn.Write([]byte(domain + "\r\n")); err != nil {
		return "", err
	}
	var output strings.Builder
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		output.WriteString(scanner.Text())
		output.WriteByte('\n')
	}
	return output.String(), scanner.Err()
}

func parse(domain, registry, registrar string) *models.WhoisResponse {
	all := registry + "\n" + registrar
	r := &models.WhoisResponse{Domain: domain}
	r.RegisteredOn = formatDate(firstValue(all, "registered on", "creation date", "created date", "created on", "created on..............", "domain registration date"))
	r.ExpiresOn = formatDate(firstValue(all, "expires on", "expires on..............", "registry expiry date", "expiration date", "expiry date", "paid-till"))
	r.UpdatedOn = formatDate(firstValue(all, "updated on", "updated date", "last updated", "last modified"))
	r.Status = valuesFor(all, "domain status", "status")
	r.NameServers = valuesFor(all, "name server", "nserver")
	if len(r.NameServers) == 0 {
		r.NameServers = trabisNameServers(all)
	}
	for index, value := range r.Status {
		r.Status[index] = formatStatus(value)
	}
	for index, value := range r.NameServers {
		r.NameServers[index] = strings.ToLower(value)
	}
	r.Registrar = models.WhoisRegistrar{Name: firstValue(all, "registrar", "sponsoring registrar"), IANAID: firstValue(all, "registrar iana id", "iana id"), Email: firstValue(all, "registrar email", "registrar url", "registrar contact email"), AbuseEmail: firstValue(all, "registrar abuse contact email", "abuse email", "abuse contact email"), AbusePhone: firstValue(all, "registrar abuse contact phone", "abuse phone", "abuse contact phone")}
	r.Registrant = contact(all, []string{"registrant", "owner"})
	r.Technical = contact(all, []string{"tech", "technical"})
	return r
}

func contact(raw string, prefixes []string) models.WhoisContact {
	return models.WhoisContact{Name: contactValue(raw, prefixes, "name"), Organization: contactValue(raw, prefixes, "organization", "org"), Street: contactValue(raw, prefixes, "street", "address"), City: contactValue(raw, prefixes, "city"), State: contactValue(raw, prefixes, "state/province", "state", "province"), PostalCode: contactValue(raw, prefixes, "postal code", "postalcode", "zip"), Country: contactValue(raw, prefixes, "country"), Phone: contactValue(raw, prefixes, "phone"), Fax: contactValue(raw, prefixes, "fax"), Email: contactValue(raw, prefixes, "email")}
}

func contactValue(raw string, prefixes []string, fields ...string) string {
	for _, prefix := range prefixes {
		for _, field := range fields {
			if value := findValue(raw, prefix+" "+field); value != "" {
				return value
			}
			if value := findValue(raw, prefix+"-"+field); value != "" {
				return value
			}
		}
	}
	return ""
}

func firstValue(raw string, keys ...string) string {
	for _, key := range keys {
		if value := findValue(raw, key); value != "" {
			return value
		}
	}
	return ""
}
func findValue(raw, key string) string {
	key = strings.ToLower(key)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && strings.EqualFold(strings.TrimSpace(parts[0]), key) {
			return cleanValue(parts[1])
		}
	}
	return ""
}
func valuesFor(raw string, keys ...string) []string {
	set := map[string]struct{}{}
	for _, line := range strings.Split(raw, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 {
			continue
		}
		for _, key := range keys {
			if strings.EqualFold(strings.TrimSpace(parts[0]), key) {
				if value := cleanValue(parts[1]); value != "" {
					set[value] = struct{}{}
				}
				break
			}
		}
	}
	values := make([]string, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}
func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	if index := strings.Index(value, " "); index > 0 && strings.HasPrefix(value[index:], " http") {
		return strings.TrimSpace(value[:index])
	}
	return value
}

func formatStatus(value string) string {
	parts := strings.Fields(value)
	if len(parts) > 0 {
		value = parts[0]
	}
	var result strings.Builder
	for index, char := range value {
		if index > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte(' ')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func formatDate(value string) string {
	if index := strings.Index(value, "T"); index != -1 {
		return value[:index]
	}
	return value
}

func trabisNameServers(raw string) []string {
	servers := []string{}
	inSection := false
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if strings.EqualFold(strings.TrimSuffix(line, ":"), "domain servers") {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		if line == "" || strings.Contains(line, ":") {
			break
		}
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.Contains(fields[0], ".") {
			servers = append(servers, strings.ToLower(fields[0]))
		}
	}
	return servers
}
