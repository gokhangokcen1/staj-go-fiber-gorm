package scanner

import (
	"os/exec"
	"strings"
)

func PingAt(ip string) bool {

	cmd := exec.Command("ping", "-n", "1", "-w", "500", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	cikti := string(output)
	return strings.Contains(cikti, "TTL=")
}
