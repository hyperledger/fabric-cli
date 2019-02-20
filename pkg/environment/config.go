/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// DefaultConfigFilename is the config filename
const DefaultConfigFilename = "config.yaml"

// Config defines the required actions for managing the config file
type Config interface {
	FromFile() (*Settings, error)
	Save() error
}

// DefaultConfig contains the default operations for managing the config file
type DefaultConfig struct {
	Filename string
	Settings *Settings
}

// FromFile retruns an instance of settings only based on config file
// This is used to manage the config file without saving overrides
func (c *DefaultConfig) FromFile() (*Settings, error) {
	v := viper.New()

	v.Set("home", c.Settings.Home.String())

	// Load config files if one exists
	v.SetConfigFile(Home(v.GetString("home")).Path(c.Filename))
	if _, err := os.Stat(v.ConfigFileUsed()); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	s := &Settings{}

	s.Config = &DefaultConfig{
		Filename: DefaultConfigFilename,
		Settings: s,
	}

	if err := v.Unmarshal(s); err != nil {
		return nil, err
	}

	return s, nil
}

// Save writes the current Settings to the config file
// Be careful not to save environment variable overrides
func (c *DefaultConfig) Save() error {
	data, err := yaml.Marshal(c.Settings)
	if err != nil {
		return err
	}

	if err := c.Settings.Home.Init(); err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.Settings.Home.Path(c.Filename), data, 0600); err != nil {
		return err
	}

	return nil
}
