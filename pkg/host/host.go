package host

import (
	"fmt"
	"math"
	"net"
	"strconv"
)

func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}

	return 0, false
}

func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil && lis == nil { //hostPort失败，并且没有listener
		return "", err
	}
	if lis != nil {
		p, ok := Port(lis) //获取listener实际监听地址
		if !ok {
			return "", fmt.Errorf("failed to extract port: %v", lis.Addr())
		}
		port = strconv.Itoa(p)
	}
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}

	addr, err = GetLocalIp()
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(addr, port), nil
}

func GetLocalIp() (ip string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	ips := make([]net.IP, 0)
	minIndex := math.MaxInt
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index >= minIndex && len(ips) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for i, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if ip.IsLoopback() {
				continue
			}
			if !ip.IsGlobalUnicast() {
				continue
			}

			if ip.IsInterfaceLocalMulticast() {
				continue
			}

			minIndex = iface.Index
			if i == 0 {
				ips = make([]net.IP, 0, 1)
			}
			ips = append(ips, ip)
			if ip.To4() != nil {
				break
			}
		}
	}

	if len(ips) != 0 {
		return ips[len(ips)-1].String(), nil
	}

	return "", nil
}
