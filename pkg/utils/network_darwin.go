//go:build darwin

package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

const (
	defaultTunnelInterface = "utun8"
)

func DefaultTunnelDevOS() string {
	return defaultTunnelInterface
}

// RouteExistsOS checks to see if a route exists for the specified prefix
func RouteExistsOS(prefix string) (bool, error) {
	if err := ValidateCIDR(prefix); err != nil {
		return false, err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return true, err
	}
	defer r.Close()
	defer w.Close()
	ns := exec.Command("netstat", "-r", "-n")
	ns.Stdout = w
	if err = ns.Start(); err != nil {
		return true, err
	}
	defer func() {
		_ = ns.Wait()
	}()

	awk := exec.Command("awk", "-v", fmt.Sprintf("ip=%s", prefix), "$1 == ip {print $1}")
	awk.Stdin = r
	var output bytes.Buffer
	awk.Stdout = &output

	// Validate the IP we're expecting is in the output
	return strings.Contains(output.String(), prefix), nil
}

// AddRoute adds a route to the specified interface
func AddRoute(prefix, dev string) error {
	_, err := RunCommand("route", "-q", "-n", "add", "-inet", prefix, "-interface", dev)
	if err != nil {
		return fmt.Errorf("v4 route add failed: %w", err)
	}

	return nil
}

// DeleteRoute deletes a darwin route for an ipv4 prefix
func DeleteRoute(prefix, dev string) error {
	_, err := RunCommand("route", "-q", "-n", "delete", "-inet", prefix, "-interface", dev)
	if err != nil {
		return fmt.Errorf("no route deleted: %w", err)
	}

	return nil
}

func DeleteInterface(logger *zap.SugaredLogger, networkInterface string) error {
	// noop
	return nil
}

func ConfigureRoutingForTunnelInterface(logger *zap.SugaredLogger, tunnelInterface string) error {
	// noop

	return nil
}
