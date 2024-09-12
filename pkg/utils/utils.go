package utils

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

var (
	ErrInvalidIpPort  = errors.New("invalid ip port mapping, please provide format <ipv4>:port")
	ErrInvalidIpv4    = errors.New("invalid ip address, please provide a valid IPv4 address")
	ErrInvalidPort    = errors.New("invalid port, please provide a valid port")
	ErrPortOutOfRange = errors.New("invalid port number, please provide between 1 - 65535")
)

// UnmarshalFlags binds flags registered in the Cobra command to a struct. The struct must have a tag of flag:<name>
// The tag must be the name of the flag. No error is returned if the flag tags are not specified.
func UnmarshalFlags(cmd *cobra.Command, out interface{}, coreValidators ...*CoreValidator) error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	if out == nil {
		return fmt.Errorf("out is nil")
	}

	t := reflect.ValueOf(out)

	// check to see if out is a pointer pointing to a struct
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expecting ptr to struct, got %T instead", out)
	}

	e := t.Elem()

	// iterate over out fields
	for i := 0; i < e.NumField(); i++ {
		// Get the structure field
		field := e.Type().Field(i)

		// Get the field name
		name := field.Name

		// Get the 'flag' tag from the field
		flag, ok := field.Tag.Lookup("flag")
		if !ok {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(""):

			val, err := cmd.Flags().GetString(flag)
			if err != nil {
				// the flag is not set is just continue
				continue
			}

			e.FieldByName(name).SetString(val)

		case reflect.TypeOf(int(0)):
			val, err := cmd.Flags().GetInt(flag)
			if err != nil {
				// the flag is not set is just continue
				continue
			}
			e.FieldByName(name).SetInt(int64(val))
		case reflect.TypeOf(true):
			val, err := cmd.Flags().GetBool(flag)
			if err != nil {
				// the flag is not set is just continue
				continue
			}
			e.FieldByName(name).SetBool(val)
		case reflect.TypeOf([]string{}):
			val, err := cmd.Flags().GetStringSlice(flag)
			if err != nil {
				// the flag is not set is just continue
				continue
			}
			e.FieldByName(name).Set(reflect.ValueOf(val))
		default:
			return fmt.Errorf("unsupported flag type: %s", field.Type)
		}
	}

	// TODO: Return user friendly error messages
	// finally validate the struct
	// Check if any validators are defined, if not create a new one and validate the struct
	if len(coreValidators) == 0 {
		cv, err := NewCoreValidator()
		if err != nil {
			return err
		}
		return cv.Validate(out)
	}

	// use provided validators
	for _, cv := range coreValidators {
		if err := cv.Validate(out); err != nil {
			return err
		}
	}

	return nil
}

// ValidateAndUnmarshal validates the cobra flag inputs and unmarshal cobra flags into the given struct.
func ValidateAndUnmarshal(cmd *cobra.Command, out interface{}, validationFns map[string]func(validator.FieldLevel) bool) error {
	cv, err := NewCoreValidator()

	if err != nil {
		return err
	}

	if validationFns != nil {
		for key, fn := range validationFns {
			// Register custom validation function
			if err := cv.RegisterValidator(key, fn); err != nil {
				return fmt.Errorf("failed to register validator: %v", err)
			}
		}
	}

	if err := UnmarshalFlags(cmd, out, cv); err != nil {
		return err
	}
	return nil
}

// GetCharCodeSum returns the sum of all charcodes for a given string
func GetCharCodeSum(str string) int {
	i := 0
	for _, val := range strings.Split(str, "") {
		i += int(val[0])
	}
	return i
}

// IsYaml returns true if the given string has a .yaml or .yml extension
func IsYaml(s string) bool {
	ext := filepath.Ext(s)
	return ext == ".yaml" || ext == ".yml"
}

// RunCommand runs the cmd and returns the combined stdout and stderr
func RunCommand(cmd ...string) (string, error) {
	output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to run %q: %w (%s)", strings.Join(cmd, " "), err, output)
	}
	return string(output), nil
}

// IsCommandAvailable checks to see if a binary is available in the current path
func IsCommandAvailable(name string) bool {
	if _, err := exec.LookPath(name); err != nil {
		return false
	}
	return true
}

// ValidateIp ensures a valid IP4/IP6 address is provided
func ValidateIp(ip string) error {
	if ip := net.ParseIP(ip); ip != nil {
		return nil
	}
	return fmt.Errorf("%s is not a valid v4 or v6 IP", ip)
}

// ValidateCIDR ensures a valid IP4/IP6 prefix is provided
func ValidateCIDR(cidr string) error {
	_, netAddr, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid v4 or v6 IP prefix: %w", cidr, err)
	}

	if cidr != netAddr.String() {
		return fmt.Errorf("Invalid network prefix provided %s, try using %s\n", cidr, netAddr.String())
	}

	return nil
}

type Ipv4Port struct {
	Ip   string
	Port int
}

// ParseToIpv4Port parses the string expecting the format "192.168.0.1:port"
// and returns a struct that contains the separated IP address and port
func ParseToIpv4Port(ipPort string) (*Ipv4Port, error) {
	// Split the input string into IP address and port using colon as the separator
	parts := strings.Split(ipPort, ":")
	if len(parts) != 2 {
		return nil, ErrInvalidIpPort
	}

	// Validate the IP address
	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, ErrInvalidIpv4
	}

	// Validate the port
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, ErrInvalidPort
	}
	if port < 1 || port > 65535 {
		return nil, ErrPortOutOfRange
	}

	return &Ipv4Port{
		Ip:   ip.To4().String(),
		Port: port,
	}, nil
}

// GetHostname returns the hostname of the device
func GetHostname() (string, error) {
	return os.Hostname()
}

// Determines if both IPs are on the same network
func IsSameNetwork(ip1 string, ip2 string, maskLength int) (bool, error) {
	if err := ValidateIp(ip1); err != nil {
		return false, err
	}
	if err := ValidateIp(ip2); err != nil {
		return false, err
	}

	// Parse the IP addresses
	addr1 := net.ParseIP(ip1)

	addr2 := net.ParseIP(ip2)

	// Calculate the network masks based on the provided mask length
	mask := net.CIDRMask(maskLength, 32) // Assuming IPv4 addresses

	// Apply the masks to the IP addresses to get the network portions
	network1 := addr1.Mask(mask)
	network2 := addr2.Mask(mask)

	// Compare the network portions to check if they are the same
	return network1.Equal(network2), nil
}

// GenerateRandomIPInRange generates a random IP address within the given subnet range
func GenerateRandomIPInRange(subnet *net.IPNet) net.IP {
	start := ipToUint(subnet.IP)
	mask := binaryToUint32(subnet.Mask)
	end := start | (^mask) // Calculate the end of the range

	randomIPValue := start + uint32(rand.Intn(int(end-start+1)))
	randomIP := uintToIP(randomIPValue)
	return randomIP
}

func ipToUint(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func binaryToUint32(mask net.IPMask) uint32 {
	ipv4Mask := mask[len(mask)-4:]
	return (uint32(ipv4Mask[0]) << 24) | (uint32(ipv4Mask[1]) << 16) | (uint32(ipv4Mask[2]) << 8) | uint32(ipv4Mask[3])
}

func uintToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func GeneratePortNumber() int {
	minPort := 1024  // Minimum usable port number
	maxPort := 65536 // Maximum usable port number (exclusive)
	return rand.Intn(maxPort-minPort+1) + minPort
}

func GenerateRandomNumber(length int) (int, error) {
	if length < 1 || length > 10 {
		return 0, errors.New("invalid number length, please provide between 1 - 10")
	}

	min := int(math.Pow10(length - 1))
	max := int(math.Pow10(length)) - 1
	return rand.Intn(max-min+1) + min, nil
}
