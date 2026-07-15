package subnet

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func StrToInt(stringDilimi []string) []int {
	intDilimi := make([]int, 0, len(stringDilimi))
	for _, str := range stringDilimi {
		sayi, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Donusum hatasi:", err)
			continue
		}
		intDilimi = append(intDilimi, sayi)
	}
	return intDilimi
}

func DecToBinary(ipBolunmus []string) []int {
	ipInt := StrToInt(ipBolunmus)
	binary := make([]int, 0, 32)
	for _, sayi := range ipInt {
		for i := 7; i >= 0; i-- {
			if (sayi & (1 << i)) != 0 {
				binary = append(binary, 1)
			} else {
				binary = append(binary, 0)
			}
		}
	}
	return binary
}

func FindSubnetBinary(subnetCIDR int) []int {
	var subnetB []int
	for i := 0; i < 32; i++ {
		if i < subnetCIDR {
			subnetB = append(subnetB, 1)
		} else {
			subnetB = append(subnetB, 0)
		}
	}
	return subnetB
}

func NetworkAddress(binaryIP []int, subnetCIDR int) []int {
	nAddress := make([]int, 0, 32)
	for i := 0; i < 32; i++ {
		if i < subnetCIDR {
			nAddress = append(nAddress, binaryIP[i])
		} else {
			nAddress = append(nAddress, 0)
		}
	}
	return nAddress
}

func BroadcastAddress(binaryIP []int, subnetCIDR int) []int {
	bAddress := make([]int, 0, 32)
	for i := 0; i < 32; i++ {
		if i < subnetCIDR {
			bAddress = append(bAddress, binaryIP[i])
		} else {
			bAddress = append(bAddress, 1)
		}
	}
	return bAddress
}

func IlkKullanilabilir(network []int) []int {
	first := make([]int, len(network))
	copy(first, network)
	for i := len(first) - 1; i >= 0; i-- {
		if first[i] == 0 {
			first[i] = 1
			break
		}
		first[i] = 0
	}
	return first
}

func SonKullanilabilir(broadcast []int) []int {
	last := make([]int, len(broadcast))
	copy(last, broadcast)
	for i := len(last) - 1; i >= 0; i-- {
		if last[i] == 1 {
			last[i] = 0
			break
		}
		last[i] = 1
	}
	return last
}

func BinaryToIP(binary []int) string {
	if len(binary) != 32 {
		return ""
	}
	var octets []string
	for i := 0; i < 32; i += 8 {
		value := 0
		for j := 0; j < 8; j++ {
			value = (value << 1) | binary[i+j]
		}
		octets = append(octets, strconv.Itoa(value))
	}
	return strings.Join(octets, ".")
}

func HostSayisi(subnetCIDR int) float64 {
	h := math.Pow(2, (float64(32-subnetCIDR))) - 2
	return h
}

func BinaryToString(binary []int, cidr int) string {
	var result strings.Builder

	for i, bit := range binary {

		if i == cidr {
			result.WriteString("|")
		}

		result.WriteString(strconv.Itoa(bit))

		if (i+1)%8 == 0 && i != 31 {
			result.WriteString("\t ")
		}
	}

	return result.String()
}

func WildcardMask(subnetMask []int) []int {
	wildcard := make([]int, len(subnetMask))

	for i := range subnetMask {
		if subnetMask[i] == 1 {
			wildcard[i] = 0
		} else {
			wildcard[i] = 1
		}
	}

	return wildcard
}

func GetIPClass(ip string) string {
	ilk, err := strconv.Atoi(strings.Split(ip, ".")[0])
	if err != nil {
		return "Unknown"
	}

	switch {
	case ilk >= 1 && ilk <= 127:
		return "Class A"
	case ilk >= 128 && ilk <= 191:
		return "Class B"
	case ilk >= 192 && ilk <= 223:
		return "Class C"
	case ilk >= 224 && ilk <= 239:
		return "Class D (Multicast)"
	case ilk >= 240 && ilk <= 255:
		return "Class E (Experimental)"
	default:
		return "Unknown"
	}
}

// binaryToUint32, 32 bitlik binary diziyi tek bir sayiya cevirir.
// Bu, IP'ler arasinda "bir sonraki IP" seklinde artirma yapabilmemiz icin gerekli.
func binaryToUint32(binary []int) uint32 {
	var result uint32
	for _, bit := range binary {
		result = (result << 1) | uint32(bit)
	}
	return result
}

// uint32ToIP, bir sayiyi tekrar "x.x.x.x" formatina cevirir.
func uint32ToIP(value uint32) string {
	return strconv.Itoa(int(value>>24&0xFF)) + "." +
		strconv.Itoa(int(value>>16&0xFF)) + "." +
		strconv.Itoa(int(value>>8&0xFF)) + "." +
		strconv.Itoa(int(value&0xFF))
}

// HostIPListesi, network adresi ile broadcast adresi ARASINDAKI
// (ikisi haric) tum host IP'lerini bir dizi olarak dondurur.
func HostIPListesi(ip string, cidr int) []string {
	ipBolunmus := strings.Split(ip, ".")
	binaryIP := DecToBinary(ipBolunmus)

	netAddr := binaryToUint32(NetworkAddress(binaryIP, cidr))
	broadAddr := binaryToUint32(BroadcastAddress(binaryIP, cidr))

	var hostlar []string
	for i := netAddr + 1; i < broadAddr; i++ {
		hostlar = append(hostlar, uint32ToIP(i))
	}
	return hostlar
}
