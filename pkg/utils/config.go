package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfiguratorProps is a struct that contains all the properties needed to create an instance of a Configurator.
type ConfiguratorProps struct {
	// The base name of the config file, without the file extension. The default value is "config" and is used
	// if the property is not set.
	DefaultFilename string

	// The file extension of the config file. The default value is "yaml" and is used if the property is not set.
	ConfigType string

	// A list of paths where the configurator should look for the config file. The default value is "." and is used
	// if the property is not set.
	ConfigPaths []string

	// When binding flags to environment variables expect that the environment variables are prefixed,
	// e.g. if the envPrefix is FOO and we have a flag like --number then number is bound
	// to FOO_NUMBER. This helps avoid conflicts.
	EnvPrefix string

	// An instance of a cobra.Command that will be used to bind flags to environment variables.
	Cmd *cobra.Command
}

// Configurator is used to read and write configuration files. It supports YAML, JSON, and TOML formats and is
// powered by Viper.
type Configurator struct {
	viper *viper.Viper
	props *ConfiguratorProps
}

// NewConfigurator creates a new instance of a Configurator with the given properties. In addition, it will
// also bind flags to environment variables and applies configuration overrides if the values are not set.
func NewConfigurator(props *ConfiguratorProps) (*Configurator, error) {

	if props == nil {
		return nil, fmt.Errorf("props cannot be nil")
	}

	if props.Cmd == nil {
		return nil, fmt.Errorf("props.Cmd cannot be nil")
	}

	c := &Configurator{
		viper: viper.New(),
		props: props,
	}

	c.viper.SetConfigName(props.DefaultFilename)

	c.viper.SetConfigType(props.ConfigType)

	// Add current folder to config search paths
	c.viper.AddConfigPath(".")

	for _, path := range props.ConfigPaths {
		c.viper.AddConfigPath(path)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := c.viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	c.viper.SetEnvPrefix(props.EnvPrefix)
	c.viper.AutomaticEnv()

	c.bindFlags()

	return c, nil
}

// bindFlags binds the flags to the environment variables and applies configuration overrides if the values
// are not set.
func (c *Configurator) bindFlags() {
	cmd := c.props.Cmd

	// Bind environment variables to flags. This sets environment variables that start with the prefix
	// but not declared anywhere else to be used by the configurator.
	envs := os.Environ()
	for _, env := range envs {
		// if the environment variable contains the prefix provided in the config props, then
		// bind the value to a key.

		// ignore all environment variables that don't contain the prefix
		if !strings.HasPrefix(env, c.props.EnvPrefix) {
			continue
		}

		// split the environment variable into key and value
		parts := strings.SplitN(env, "=", 2)

		// if the environment variable doesn't contain a value, then ignore it
		if len(parts) < 2 {
			continue
		}

		// replace the prefix with an empty string
		key := strings.Replace(parts[0], fmt.Sprintf("%s_", c.props.EnvPrefix), "", 1)

		// replace "_" with "-"
		key = strings.ToLower(strings.ReplaceAll(key, "_", "-"))

		// set the value of the environment variable override it
		c.viper.Set(key, parts[1])

	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {

		// Start by binding cmd.Flag to viper
		_ = c.viper.BindPFlag(flag.Name, flag)

		// Replace all instances of "-" with "_" in the flag name and bind it to the environment variable.
		envVarSuffix := strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))

		if strings.Contains(flag.Name, "-") {
			_ = c.viper.BindEnv(flag.Name, fmt.Sprintf("%s_%s", c.props.EnvPrefix, envVarSuffix))

		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !flag.Changed && c.viper.IsSet(flag.Name) {
			val := c.viper.Get(flag.Name)
			_ = cmd.Flags().Set(flag.Name, fmt.Sprintf("%v", val))
		}

		// However, if the flag is set, then update viper with the new value
		if flag.Changed {
			val := flag.Value.String()
			c.viper.Set(flag.Name, val)
		}

	})
}

// Viper returns the viper instance used by the configurator.
func (c *Configurator) Viper() *viper.Viper {
	return c.viper
}

// Unmarshal validates and unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func (c *Configurator) Unmarshal(out interface{}) error {
	v := reflect.ValueOf(out)

	// check to see if out is of kind struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expecting ptr to struct, got %T instead", out)
	}

	if err := c.viper.Unmarshal(out); err != nil {
		return errors.New(fmt.Sprintf("unable to unmarshal config: %s", err.Error()))
	}

	// validate the config
	cv, err := NewCoreValidator()
	if err != nil {
		return err
	}
	return cv.Validate(out)
}
