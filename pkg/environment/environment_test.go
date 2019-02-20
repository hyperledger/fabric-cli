/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSettings(t *testing.T) {
	tests := []struct {
		// test name
		name string

		// environemnt variables
		env map[string]string

		// expected output
		home    Home
		plugins string
	}{
		{
			name:    "Environment Variables Not Set",
			home:    DefaultHome,
			plugins: DefaultHome.Plugins(),
		},
		{
			name: "Environment Variables Set",
			env: map[string]string{
				"FABRIC_HOME": "/foo",
			},
			home:    Home("/foo"),
			plugins: "/foo/plugins",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				os.Setenv(k, v)
			}

			settings, err := GetSettings()

			assert.Nil(t, err)
			assert.Equal(t, settings.Home, test.home)
			assert.Equal(t, settings.Home.Plugins(), test.plugins)

			for k := range test.env {
				os.Unsetenv(k)
			}
		})
	}
}

func TestSettingsPluginEnv(t *testing.T) {
	settings, err := GetSettings()

	assert.Nil(t, err)
	assert.NotNil(t, settings)
	assert.Zero(t, os.Getenv("FABRIC_HOME"))

	settings.SetupPluginEnv()

	assert.Equal(t, os.Getenv("FABRIC_HOME"), settings.Home.String())
}
