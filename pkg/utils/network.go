package utils

import (
	"net"

	"go.uber.org/zap"
)

// GetInterfacesWithLocalAddr get all interfaces with local ip addresses except for the tunnel interface
func GetInterfacesWithLocalAddr() ([]net.Interface, error) {
	var result []net.Interface

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, networkInterface := range interfaces {
		addresses, err := networkInterface.Addrs()
		if err != nil {
			return nil, err
		}

		// Validate the interface has a local address
		for _, addr := range addresses {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				result = append(result, networkInterface)
			}
		}
	}

	return result, nil
}

// ParseIPNet return an IPNet from a string
func ParseIPNet(s string) (*net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: ipNet.Mask}, nil
}

// IfaceExists returns true if the input matches a net interface
func IfaceExists(logger *zap.SugaredLogger, networkInterface string) bool {
	_, err := net.InterfaceByName(networkInterface)
	if err != nil {
		logger.Debugf("existing link not found: %s", networkInterface)
		return false
	}

	return true
}
