/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Settings contains the combination of all environment settings
type Settings struct {
	Config `yaml:"-"`

	Home           Home       `yaml:"-"`
	Streams        Streams    `yaml:"-"`
	ActiveProfile  string     `mapstructure:"active_profile" yaml:"active_profile"`
	Profiles       []*Profile `mapstructure:"profiles" yaml:"profiles"`
	DisablePlugins bool       `mapstructure:"disable_plugins" yaml:"disable_plugins"`
}

// SetupPluginEnv sets the environment variables that are important to plugins
// This is needed because environment variables are not populated with defaults
func (s *Settings) SetupPluginEnv() {
	for k, v := range map[string]string{
		"FABRIC_HOME": s.Home.String(),
	} {
		os.Setenv(k, v) // nolint: errcheck
	}
}

// NewDefaultSettings returns an instance of settings with default config
// implementation
func NewDefaultSettings() *Settings {
	s := &Settings{}
	s.Config = &DefaultConfig{
		Filename: DefaultConfigFilename,
		Settings: s,
	}

	return s
}

// GetSettings populates a Settings struct based on viper precendence
// Highest precedence to lowest:
// 		Env > Config File > Defaults
// This can support flags in the future
func GetSettings() (*Settings, error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values(establish field type)
	v.SetDefault("home", DefaultHome)
	v.SetDefault("streams", DefaultStreams)
	v.SetDefault("profiles", []*Profile{})
	v.SetDefault("disable_plugins", false)

	// Load environment variable overrides
	v.SetEnvPrefix("fabric")
	v.BindEnv("home")            // nolint: errcheck
	v.BindEnv("disable_plugins") // nolint: errcheck

	// Load config files if one exists
	v.SetConfigFile(Home(v.GetString("home")).Path(DefaultConfigFilename))
	if _, err := os.Stat(v.ConfigFileUsed()); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	s := NewDefaultSettings()

	if err := v.Unmarshal(s); err != nil {
		return nil, err
	}

	return s, nil
}
