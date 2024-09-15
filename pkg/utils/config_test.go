package utils

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var config = `
foo: barYaml
baz: quxYaml
baz-baz: quxQuxYaml
`

type ConfigTestSuite struct {
	suite.Suite
	props            *ConfiguratorProps
	tmpConfigFileDir string
}

func (suite *ConfigTestSuite) SetupSuite() {
	dir, err := ioutil.TempDir("", "configurator_test")
	require.Nil(suite.T(), err)

	suite.tmpConfigFileDir = dir

	filename := dir + "/config.yaml"

	err = ioutil.WriteFile(filename, []byte(config), 0644)
	require.Nil(suite.T(), err)

}

func (suite *ConfigTestSuite) SetupTest() {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	var foo string
	cmd.Flags().StringVarP(&foo, "foo", "f", "barCmd", "Test")

	suite.props = &ConfiguratorProps{
		DefaultFilename: "config",
		ConfigType:      "",
		ConfigPaths:     []string{suite.tmpConfigFileDir},
		EnvPrefix:       "PC_TEST",
		Cmd:             cmd,
	}
}

func (suite *ConfigTestSuite) TearDownSuite() {
	_ = os.RemoveAll(suite.tmpConfigFileDir)
}

func (suite *ConfigTestSuite) TearDownTest() {

	envs := os.Environ()
	for _, env := range envs {
		if strings.HasPrefix(env, "PC_TEST") {
			// split the environment variable into key and value
			parts := strings.SplitN(env, "=", 2)

			// if the environment variable doesn't contain a value, then ignore it
			if len(parts) < 2 {
				continue
			}
			_ = os.Unsetenv(parts[0])
		}
	}
}

func (suite *ConfigTestSuite) TestNewConfiguratorNilProps() {
	_, err := NewConfigurator(nil)
	require.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "props cannot be nil", err.Error())
}

func (suite *ConfigTestSuite) TestNewConfiguratorProps() {
	props := &ConfiguratorProps{
		DefaultFilename: "",
		ConfigType:      "",
		ConfigPaths:     nil,
		EnvPrefix:       "",
		Cmd:             nil,
	}

	_, err := NewConfigurator(props)
	require.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "props.Cmd cannot be nil", err.Error())

	props.Cmd = &cobra.Command{}
	configurator, err := NewConfigurator(props)

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), configurator)
}

func (suite *ConfigTestSuite) TestNewConfiguratorFile() {
	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barYaml", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
}

// TestNewConfigruatorEnv1 tests that a value defined in the environment variable overrides
// value set in a file.
func (suite *ConfigTestSuite) TestNewConfiguratorEnv1() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")
	_ = os.Setenv("PC_TEST_BAZ_BAZ", "quxQuxEnv")

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barEnv", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
	assert.Equal(suite.T(), "quxQuxEnv", configurator.Viper().Get("baz-baz"))
}

func (suite *ConfigTestSuite) TestNewConfiguratorEnv2() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")
	_ = os.Setenv("PC_TEST_BAZ_BAZ_BAZ", "quxQuxQuxEnv")
	_ = os.Setenv("PC_TEST_QUX", "quxEnv")

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barEnv", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
	assert.Equal(suite.T(), "quxQuxYaml", configurator.Viper().Get("baz-baz"))
	assert.Equal(suite.T(), "quxQuxQuxEnv", configurator.Viper().Get("baz-baz-baz"))
	assert.Equal(suite.T(), "quxEnv", configurator.Viper().Get("qux"))
}

func (suite *ConfigTestSuite) TestNewConfiguratorCmd1() {

	cmd := suite.props.Cmd
	cmd.SetArgs([]string{"--foo", "barCli"})
	_ = cmd.Execute()
	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barCli", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
}

func (suite *ConfigTestSuite) TestNewConfiguratorCmd2() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barEnv", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
}

func (suite *ConfigTestSuite) TestNewConfiguratorCmd3() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")

	var foo string
	cmd := suite.props.Cmd
	cmd.Flags().StringVarP(&foo, "fooCmd", "k", "barCmd", "Test")

	err := cmd.Execute()
	assert.Nil(suite.T(), err)

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barEnv", configurator.Viper().Get("foo"))
	assert.Equal(suite.T(), "barCmd", configurator.Viper().Get("fooCmd"))
	assert.Equal(suite.T(), "quxYaml", configurator.Viper().Get("baz"))
}

func (suite *ConfigTestSuite) TestNewConfiguratorUnmarshal() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")
	var foo string
	cmd := suite.props.Cmd
	cmd.Flags().StringVarP(&foo, "fooCmd", "k", "barCmd", "Test")

	err := cmd.Execute()
	assert.Nil(suite.T(), err)

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	c := struct {
		Foo    string `mapstructure:"foo"`
		Baz    string `mapstructure:"baz"`
		BazBaz string `mapstructure:"baz-baz"`
	}{}

	err = configurator.Unmarshal(&c)
	require.Nil(suite.T(), err)

	assert.Equal(suite.T(), "barEnv", c.Foo)
	assert.Equal(suite.T(), "quxYaml", c.Baz)
	assert.Equal(suite.T(), "quxQuxYaml", c.BazBaz)
}

// TestConfiguratorValidator tests validation happy path
func (suite *ConfigTestSuite) TestConfiguratorValidator1() {

	// When an env variable is set, it should override the config file
	_ = os.Setenv("PC_TEST_FOO", "barEnv")

	var foo string
	cmd := suite.props.Cmd
	cmd.Flags().StringVarP(&foo, "fooCmd", "k", "barCmd", "Test")

	err := cmd.Execute()
	assert.Nil(suite.T(), err)

	configurator, err := NewConfigurator(suite.props)

	require.Nil(suite.T(), err)

	c := struct {
		Foo    string `mapstructure:"foo" validate:"endswith=Env"`
		Baz    string `mapstructure:"baz"`
		BazBaz string `mapstructure:"baz-baz"`
	}{}

	err = configurator.Unmarshal(&c)

	require.Nil(suite.T(), err)
}

// TestConfiguratorValidator2 tests validation error path
func (suite *ConfigTestSuite) TestConfiguratorValidator2() {
	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	c := struct {
		Foo    string `mapstructure:"foo" validate:"endswith=Env"`
		Baz    string `mapstructure:"baz"`
		BazBaz string `mapstructure:"baz-baz"`
	}{}

	err = configurator.Unmarshal(&c)
	require.NotNil(suite.T(), err)

	assert.Equal(suite.T(),
		"Key: 'Foo' Error:Field validation for 'Foo' failed on the 'endswith' tag: failed constraint endswith=Env, received: barYaml",
		err.Error())
}

// TestConfiguratorValidator3 tests all other execution paths
func (suite *ConfigTestSuite) TestConfiguratorValidator3() {

	configurator, err := NewConfigurator(suite.props)
	require.Nil(suite.T(), err)

	assert.NotNil(suite.T(), configurator.Viper())

	var s string
	err = configurator.Unmarshal(&s)
	assert.Equal(suite.T(), "expecting ptr to struct, got *string instead", err.Error())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
