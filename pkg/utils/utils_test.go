package utils_test

import (
	_ "embed"
	"net"
	"strings"
	"testing"

	"yuka/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Config   string   `flag:"config"`
	Env      string   `flag:"env"`
	Replicas int      `flag:"replicas"`
	Tags     []string `flag:"tags"`
}

type testStructNoFlag struct {
	Config string
	Env    string
}

type testStructPartialFlag struct {
	Config string `flag:"config"`
	Env    string
}

type testStructUnsupportedType struct {
	Replicas float64 `flag:"replicas"`
}

type testStructValidate struct {
	Config string `flag:"config" validate:"required"`
	Env    string `flag:"env" validate:"required"`
}

type testValidateUnmarshal struct {
	Config  string `flag:"config" validate:"required"`
	Env     string `flag:"env" validate:"required"`
	Project string `flag:"project" validate:"isPci"`
}

func isPci(fl validator.FieldLevel) bool {
	if strings.HasPrefix(fl.Field().String(), "PCI") {
		return true
	}
	return false
}

func TestUnmarshalFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")
	cmd.Flags().IntP("replicas", "r", 1, "Number of replicas")
	cmd.Flags().StringSliceP("tags", "t", []string{"testTag1", "testTag2"},
		"Tags to apply to the manifest")

	out := &testStruct{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.Nil(t, err)

	assert.Equal(t, "testConfig", out.Config)
	assert.Equal(t, "testEnv", out.Env)
	assert.Equal(t, 1, out.Replicas)
	assert.Equal(t, []string{"testTag1", "testTag2"}, out.Tags)
}

func TestUnmarshalFlagsWithNoFlagTag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")

	out := &testStructNoFlag{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.Nil(t, err)

	assert.Equal(t, "", out.Config)
	assert.Equal(t, "", out.Env)
}

func TestUnmarshalPartialFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")
	cmd.Flags().IntP("replicas", "r", 1, "Number of replicas")
	cmd.Flags().StringSliceP("tags", "t", []string{"testTag1", "testTag2"},
		"Tags to apply to the manifest")

	out := &testStructPartialFlag{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.Nil(t, err)

	assert.Equal(t, "testConfig", out.Config)
	assert.Equal(t, "", out.Env)
}

func TestUnmarshalUnsupportedFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")
	cmd.Flags().Float64P("replicas", "r", 1, "Number of replicas")

	out := &testStructUnsupportedType{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported flag type: float64", err.Error())
}

func TestUnmarshalFlagsWithValidationPass(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")

	out := &testStructValidate{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.Nil(t, err)

	assert.Equal(t, "testConfig", out.Config)
	assert.Equal(t, "testEnv", out.Env)

}

func TestUnmarshalFlagsWithValidationFail(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")

	out := &testStructValidate{}
	err := utils.UnmarshalFlags(cmd, out)
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'testStructValidate.Env' Error:Field validation for 'Env' failed on the 'required' "+
		"tag: failed constraint required=, received: ", err.Error())

}

func TestUnmarshalFlagsOutNil(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")

	err := utils.UnmarshalFlags(cmd, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "out is nil", err.Error())
}

func TestUnmarshalFlagsCmdNil(t *testing.T) {
	err := utils.UnmarshalFlags(nil, &testStruct{})
	assert.NotNil(t, err)
	assert.Equal(t, "cmd is nil", err.Error())
}

func TestValidateAndUnmarshalCmdNil(t *testing.T) {
	err := utils.ValidateAndUnmarshal(nil, &testStruct{}, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "cmd is nil", err.Error())
}

func TestValidateAndUnmarshalOptionsNil(t *testing.T) {
	err := utils.ValidateAndUnmarshal(&cobra.Command{}, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "out is nil", err.Error())
}

func TestValidateAndUnmarshalValidationsNil(t *testing.T) {
	err := utils.ValidateAndUnmarshal(&cobra.Command{}, &testStruct{}, nil)
	assert.Nil(t, err)
}

func TestValidateAndUnmarshalNoProject(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")

	out := &testValidateUnmarshal{}

	var validationFns = map[string]func(validator.FieldLevel) bool{
		"isPci": isPci,
	}

	err := utils.ValidateAndUnmarshal(cmd, out, validationFns)
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'testValidateUnmarshal.Project' Error:Field validation for "+
		"'Project' failed on the 'isPci' tag: failed constraint isPci=, received: ", err.Error())
}

func TestValidateAndUnmarshalInvalidProject(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")
	cmd.Flags().StringP("project", "p", "testProject", "Project to generate manifest for")

	out := &testValidateUnmarshal{}

	var validationFns = map[string]func(validator.FieldLevel) bool{
		"isPci": isPci,
	}

	err := utils.ValidateAndUnmarshal(cmd, out, validationFns)
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'testValidateUnmarshal.Project' Error:Field validation for 'Project' failed "+
		"on the 'isPci' tag: failed constraint isPci=, received: testProject", err.Error())
}

func TestValidateAndUnmarshal(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("config", "c", "testConfig", "Path to config file")
	cmd.Flags().StringP("env", "e", "testEnv", "Environment to generate manifest for")
	cmd.Flags().StringP("project", "p", "PCI1", "Project to generate manifest for")

	out := &testValidateUnmarshal{}

	var validationFns = map[string]func(validator.FieldLevel) bool{
		"isPci": isPci,
	}

	err := utils.ValidateAndUnmarshal(cmd, out, validationFns)
	assert.Nil(t, err)

	assert.Equal(t, "testConfig", out.Config)
	assert.Equal(t, "testEnv", out.Env)
}

func TestIsYamlFn(t *testing.T) {
	assert.True(t, utils.IsYaml("test.yaml"))
	assert.True(t, utils.IsYaml("test.yml"))
	assert.False(t, utils.IsYaml("test.json"))
}

// ParseToIpv4Port tests
func TestInvalidIpPort(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("asdojoasdj")
	assert.Nil(t, ipv4)
	assert.Equal(t, err, utils.ErrInvalidIpPort)
}

func TestInvalidIpAddress(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("10.0.0:1280")
	assert.Nil(t, ipv4)
	assert.Equal(t, err, utils.ErrInvalidIpv4)
}

func TestInvalidPort(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("192.168.0.1:blah")
	assert.Nil(t, ipv4)
	assert.Equal(t, err, utils.ErrInvalidPort)
}

func TestPortNumberTooSmall(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("192.168.0.1:0")
	assert.Nil(t, ipv4)
	assert.Equal(t, err, utils.ErrPortOutOfRange)
}

func TestPortNumberTooLarge(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("192.168.0.1:65536")
	assert.Nil(t, ipv4)
	assert.Equal(t, err, utils.ErrPortOutOfRange)
}

func TestSuccessfullyParsesIpv4Port(t *testing.T) {
	ipv4, err := utils.ParseToIpv4Port("192.168.0.1:51280")
	assert.Nil(t, err)
	assert.Equal(t, ipv4, &utils.Ipv4Port{
		Ip:   "192.168.0.1",
		Port: 51280,
	})
}

func TestIsSameNetwork(t *testing.T) {
	// Test cases
	testCases := []struct {
		name            string
		myIPAddress     string
		targetIPAddress string
		subnetMaskLen   int
		expected        bool
		shouldError     bool
	}{
		// Test cases where both IP addresses are on the same network
		{"Same network", "192.168.1.100", "192.168.1.200", 24, true, false},
		{"Same network with different IPs", "192.168.1.1", "192.168.1.2", 24, true, false},

		// Test cases where both IP addresses are NOT on the same network
		{"Different networks", "192.168.1.100", "192.168.2.200", 24, false, false},
		{"Different networks with different mask", "192.168.1.100", "192.169.1.200", 16, false, false},

		// // Test case with invalid IP address
		{"Invalid IP address", "invalid", "192.168.1.100", 24, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.IsSameNetwork(tc.myIPAddress, tc.targetIPAddress, tc.subnetMaskLen)

			if tc.shouldError {
				assert.Error(t, err, "Expected an error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got an error")
			}

			assert.Equal(t, tc.expected, result, "Mismatch in expected result")
		})
	}
}

func TestGenerateRandomIPInRange(t *testing.T) {
	subnet := &net.IPNet{
		IP:   net.IPv4(192, 168, 1, 0),
		Mask: net.IPMask(net.IPv4Mask(255, 255, 255, 0)),
	}

	for i := 0; i < 100; i++ {
		randomIP := utils.GenerateRandomIPInRange(subnet)

		if !subnet.Contains(randomIP) {
			t.Errorf("Generated IP %s is not within the subnet range %s", randomIP, subnet.String())
		}
	}

	subnet = &net.IPNet{
		IP:   net.IPv4(10, 0, 0, 0),
		Mask: net.IPMask(net.IPv4Mask(255, 0, 0, 0)),
	}

	for i := 0; i < 1000; i++ {
		randomIP := utils.GenerateRandomIPInRange(subnet)

		if !subnet.Contains(randomIP) {
			t.Errorf("Generated IP %s is not within the subnet range %s", randomIP, subnet.String())
		}
	}

}

func TestGenerateRandomNumber(t *testing.T) {
	// Test valid input
	for i := 1; i <= 10; i++ {
		num, err := utils.GenerateRandomNumber(i)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if num < int((10^(i-1))) || num > int((10^i)-1) {
			t.Errorf("Generated number %d is not within range [%d, %d]", num, int((10 ^ (i - 1))), int((10^i)-1))
		}
	}

	// Test invalid input
	_, err := utils.GenerateRandomNumber(0)
	if err == nil || err.Error() != "invalid token length, please provide between 1 - 10" {
		t.Errorf("Expected error 'invalid token length, please provide between 1 - 10', but got %v", err)
	}

	_, err = utils.GenerateRandomNumber(11)
	if err == nil || err.Error() != "invalid token length, please provide between 1 - 10" {
		t.Errorf("Expected error 'invalid token length, please provide between 1 - 10', but got %v", err)
	}
}
