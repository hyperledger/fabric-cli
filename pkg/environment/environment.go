/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"fmt"
	"os"
)

// Settings contains all of the environment settings
type Settings struct {
	Home           Home
	Streams        Streams
	DisablePlugins bool
}

// parseEnvironments overrides all default settings
func (s *Settings) parseEnvironment() {
	if v, ok := os.LookupEnv("FABRIC_HOME"); ok {
		s.Home = Home(v)
	}
	if v, ok := os.LookupEnv("FABRIC_DISABLE_PLUGINS"); ok && v == "true" {
		s.DisablePlugins = true
	}
}

// setEnvironment sets all environment variables after settings is loaded
// This is needed because default values won't be found in the environment
func (s *Settings) setEnvironment() error {
	if err := os.Setenv("FABRIC_HOME", s.Home.String()); err != nil {
		return err
	}

	if err := os.Setenv("FABRIC_DISABLE_PLUGINS", fmt.Sprint(s.DisablePlugins)); err != nil {
		return err
	}

	return nil
}

// GetSettings returns the settings for the current environment
func GetSettings() (*Settings, error) {
	s := &Settings{
		Home:    DefaultHome,
		Streams: *DefaultStreams,
	}

	s.parseEnvironment()

	if err := s.setEnvironment(); err != nil {
		return nil, err
	}

	return s, nil
}
