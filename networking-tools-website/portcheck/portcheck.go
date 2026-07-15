package portcheck

import (
	"fmt"
	"net"
	"time"
)

func KontrolEt(ip string, port int) (bool, string) {
	adres := fmt.Sprintf("%v:%v", ip, port)
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", adres, timeout)
	if err != nil {
		return false, err.Error()
	}
	defer conn.Close()

	return true, ""
}
