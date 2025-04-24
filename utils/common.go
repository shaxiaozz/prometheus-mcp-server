package utils

import (
	"github.com/wonderivan/logger"
	"net"
)

func GetIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Error("Get IP Addr Error: ", err)
		return "localhost"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// IPv4
			if ipNet.IP.To4() != nil {
				logger.Info("Get IP Addr: ", ipNet.IP.String())
				return ipNet.IP.String()
			}
		}
	}
	return "localhost"
}
