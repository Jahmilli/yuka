package stun

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/libp2p/go-reuseport"
	"github.com/pion/stun"
	"go.uber.org/zap"
)

func RequestWithReusePort(logger *zap.SugaredLogger, stunServer string, srcPort int) (netip.AddrPort, error) {
	logger.Debugf("dialing stun Server %s with srcPort %v", stunServer, srcPort)
	conn, err := reuseport.Dial("udp4", fmt.Sprintf(":%d", srcPort), stunServer)
	if err != nil {
		logger.Errorf("stun dialing timed out %v", err)
		return netip.AddrPort{}, fmt.Errorf("failed to dial stun Server %s: %w", stunServer, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c, err := stun.NewClient(conn)
	if err != nil {
		logger.Error(err)
		return netip.AddrPort{}, err
	}
	defer func() {
		_ = c.Close()
	}()

	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN Server, waiting for response message.
	var xorAddr stun.XORMappedAddress
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			if res.Error.Error() == "transaction is timed out" {
				logger.Debugf("STUN transaction timed out, if this continues check if a firewall is blocking UDP connections to %s", stunServer)
			} else {
				logger.Debug(res.Error)
			}
			return
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		if err := xorAddr.GetFrom(res.Message); err != nil {
			return
		}
	}); err != nil {
		return netip.AddrPort{}, err
	}

	xorBinding, err := netip.ParseAddrPort(xorAddr.String())
	if err != nil {
		return netip.AddrPort{}, fmt.Errorf("failed to parse a valid address:port binding from the stun response: %w", err)
	}
	logger.Debugf("reflexive binding is: %s", xorBinding.String())

	return xorBinding, nil
}

type SymmetricNatResponse struct {
	IsSymmetric          bool
	ReflexiveAddressIPv4 netip.AddrPort
}

// CheckSymmetricNat determines if the joining node is within a symmetric NAT cone
func CheckSymmetricNat(ctx context.Context, logger *zap.SugaredLogger, listenPort int) (*SymmetricNatResponse, error) {
	symmetricNat := false
	stunServer1 := NextServer()
	stunServer2 := NextServer()

	stunAddr1, err := Request(logger, stunServer1, listenPort)
	var nodeReflexiveAddressIPv4 netip.AddrPort
	if err != nil {
		return nil, err
	} else {
		nodeReflexiveAddressIPv4 = stunAddr1
	}

	isSymmetric := false
	stunAddr2, err := Request(logger, stunServer2, listenPort)
	if err != nil {
		return nil, err
	} else {
		isSymmetric = stunAddr1.String() != stunAddr2.String()
	}

	if stunAddr1.Addr().String() != "" {
		logger.Debugf("first NAT discovery STUN request returned: %s", stunAddr1.String())
	} else {
		logger.Debugf("first NAT discovery STUN request returned an empty value")
	}

	if stunAddr2.Addr().String() != "" {
		logger.Debugf("second NAT discovery STUN request returned: %s", stunAddr2.String())
	} else {
		logger.Debugf("second NAT discovery STUN request returned an empty value")
	}

	if isSymmetric {
		symmetricNat = true
		logger.Infof("Symmetric NAT is detected, this node will be provisioned in relay mode only")
	}

	return &SymmetricNatResponse{
		IsSymmetric:          symmetricNat,
		ReflexiveAddressIPv4: nodeReflexiveAddressIPv4,
	}, nil
}
