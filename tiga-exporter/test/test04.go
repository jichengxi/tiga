package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, i := range addrs {
		if ipNet, ok := i.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				if strings.HasPrefix(ipNet.IP.String(), "172.21") ||
					strings.HasPrefix(ipNet.IP.String(), "172.20") ||
					strings.HasPrefix(ipNet.IP.String(), "172.28") {
					fmt.Println(ipNet.IP.String())
				}
			}
		}
	}
}
