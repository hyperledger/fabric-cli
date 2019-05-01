/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"os"

	"github.com/spf13/pflag"
)

// Settings contains environment configuration details
type Settings struct {
	Home           Home
	Streams        Streams
	DisablePlugins bool
	Config         *Config

	configFilename string
}

// AddFlags appeneds settings flags onto an existing flag set
func (s *Settings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar((*string)(&s.Home), "home", DefaultHome.String(), "set path to configuration files")
	fs.BoolVar(&s.DisablePlugins, "disable-plugins", false, "disable plugins")
}

// Init populates the settings based on a precedence:
// 	Flag > Env > Config File > Defaults
func (s *Settings) Init(fs *pflag.FlagSet) error {
	// if the flag was not set, check the corresponding environment variable
	// must resolve home value before loading config from file
	for flag, env := range map[string]string{
		"home":            "FABRIC_HOME",
		"disable-plugins": "FABRIC_DISABLE_PLUGINS",
	} {
		if fs.Changed(flag) {
			continue
		}

		if v, ok := os.LookupEnv(env); ok {
			fs.Set(flag, v)
		}
	}

	if err := s.Home.Init(); err != nil {
		return err
	}

	// if the config file does not exist, continue since it could be a new user or home
	if err := s.Config.LoadFromFile(s.Home.Path(s.configFilename)); err != nil &&
		!os.IsNotExist(err) {
		return err
	}

	return nil
}

// ModifyConfig loads the config file and updates it based on actions
func (s *Settings) ModifyConfig(actions ...Action) error {
	if len(actions) == 0 {
		return nil
	}

	// build a new config to prevent saving overrides
	fromFile := NewConfig()

	if err := s.Home.Init(); err != nil {
		return err
	}

	// if config file does not exist, create a new one
	if err := fromFile.LoadFromFile(s.Home.Path(s.configFilename)); err != nil &&
		!os.IsNotExist(err) {
		return err
	}

	for _, action := range actions {
		action(fromFile)

		// update settings in case it continues to be used after modification
		action(s.Config)
	}

	return fromFile.Save(s.Home.Path(s.configFilename))
}

// SetupPluginEnvironment sets the environment variables that are important to plugins
// This is needed because environment variables are not populated with defaults
func (s *Settings) SetupPluginEnvironment() {
	for k, v := range map[string]string{
		"FABRIC_HOME": s.Home.String(),
	} {
		os.Setenv(k, v) // nolint: errcheck
	}
}

// NewDefaultSettings returns settings populated with default values
func NewDefaultSettings() *Settings {
	return &Settings{
		Home:    DefaultHome,
		Streams: DefaultStreams,
		Config:  NewConfig(),

		configFilename: DefaultConfigFilename,
	}
}
