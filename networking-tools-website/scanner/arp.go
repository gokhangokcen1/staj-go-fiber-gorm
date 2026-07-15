package scanner

import (
	"os/exec"
	"regexp"
	"strings"
)

var macRegex = regexp.MustCompile(`([0-9a-fA-F]{2}-){5}[0-9a-fA-F]{2}`)

func MacAdresiBul(ip string) string {
	cmd := exec.Command("arp", "-a", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	cikti := string(output)
	eslesme := macRegex.FindString(cikti)
	return strings.ToUpper(eslesme)
}
