package host

import (
	"fmt"
	"math"
	"net"
	"os"
	"testing"
)

func TestNetInterfaces(t *testing.T) {
	ifaces, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range ifaces {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String(), int(^uint(0)>>1), math.MaxInt)
			}

		}
	}

	ip, _ := GetLocalIp()
	fmt.Println("get local ip:", ip)

}

func TestExtract(t *testing.T) {
	ip, err := GetLocalIp()
	if err != nil {
		t.Error(err)
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, 9999))
	if err != nil {
		t.Error(err)
		return
	}
	defer lis.Close()

	addr, err := Extract("", lis)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("addr:", addr)

}
