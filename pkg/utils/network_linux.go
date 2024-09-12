//go:build linux

package utils

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

const (
	defaultTunnelInterface = "wg0"
)

func DefaultTunnelDevOS() string {
	return defaultTunnelInterface
}

func AddRoute(prefix, dev string) error {
	link, err := netlink.LinkByName(dev)
	if err != nil {
		return fmt.Errorf("failed to lookup netlink device %s: %w", dev, err)
	}

	destNet, err := ParseIPNet(prefix)
	if err != nil {
		return fmt.Errorf("failed to parse a valid network address from %s: %w", prefix, err)
	}

	return route(destNet, link)
}

// AddRoute adds a netlink route pointing to the linux device
func route(ipNet *net.IPNet, dev netlink.Link) error {
	return netlink.RouteAdd(&netlink.Route{
		LinkIndex: dev.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Dst:       ipNet,
	})
}

// RouteExistsOS checks netlink routes for the destination prefix
func RouteExistsOS(prefix string) (bool, error) {
	destNet, err := ParseIPNet(prefix)
	if err != nil {
		return false, fmt.Errorf("failed to parse a valid network address from %s: %w", prefix, err)
	}

	destRoute := &netlink.Route{Dst: destNet}
	family := netlink.FAMILY_V6
	if destNet.IP.To4() != nil {
		family = netlink.FAMILY_V4
	}

	match, err := netlink.RouteListFiltered(family, destRoute, netlink.RT_FILTER_DST)
	if err != nil {
		return true, fmt.Errorf("error retrieving netlink routes: %w", err)
	}

	if len(match) > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteInterface(logger *zap.SugaredLogger, networkInterface string) error {
	if IfaceExists(logger, networkInterface) {
		_, err := RunCommand("ip", "link", "del", networkInterface)
		if err != nil {
			logger.Errorf("failed to delete the ip link interface: %v", err)
			return err
		}
	}
	return nil
}

func ConfigureRoutingForTunnelInterface(logger *zap.SugaredLogger, tunnelInterface string) error {
	if _, err := RunCommand("sysctl", "net.ipv4.ip_forward=1"); err != nil {
		return err
	}

	// Enable forwarding between interfaces
	// TODO: This isn't great, it'll get interfaces like Docker, Bridge interfaces used for Kind etc which we don't want.
	// Maybe we can get the interface that matches the subnet exposed as part of the flag but need to figure it out in proper implementation.
	interfaces, err := GetInterfacesWithLocalAddr()
	if err != nil {
		return err
	}
	for _, val := range interfaces {
		if val.Name == tunnelInterface {
			// We're updating existing interfaces to forward to tunnel interface so we should skip here to prevent forwarding from
			// tunnel interface to tunnel interface
			continue
		}
		// TODO: Add in checks if we need to do this
		_, err := RunCommand("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", val.Name, "-j", "MASQUERADE")
		if err != nil {
			return err
		}
		if _, err := RunCommand("iptables", "-A", "FORWARD", "-i", tunnelInterface, "-o", val.Name, "-j", "ACCEPT"); err != nil {
			return err
		}
		if _, err := RunCommand("iptables", "-A", "FORWARD", "-i", val.Name, "-o", tunnelInterface, "-j", "ACCEPT"); err != nil {
			return err
		}
		logger.Debugf("configuring ip tables for tunnel interface %s and network interface %s", tunnelInterface, val.Name)
		if _, err := RunCommand("iptables", "-A", "FORWARD", "-i", tunnelInterface, "-o", val.Name, "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT"); err != nil {
			logger.Errorf("Failed to configure ip table forwarding from tunnel interface %s to default interface %s\n", tunnelInterface, val.Name)
			return err
		}
		if _, err := RunCommand("iptables", "-A", "FORWARD", "-i", val.Name, "-o", tunnelInterface, "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT"); err != nil {
			logger.Errorf("Failed to configure ip table forwarding from default interface %s to tunnel interface %s\n", val.Name, tunnelInterface)
			return err
		}
	}

	return nil
}
