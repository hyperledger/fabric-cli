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

func TestFromConfigFile(t *testing.T) {
	os.Setenv("FABRIC_HOME", os.TempDir())

	before, err := GetSettings()

	assert.Nil(t, err)
	assert.NotNil(t, before)

	config, err := before.FromFile()
	assert.Nil(t, err)
	assert.NotNil(t, config)

	config.ActiveProfile = "foo"

	saveErr := config.Save()
	assert.Nil(t, saveErr)

	after, err := GetSettings()

	assert.Nil(t, err)
	assert.NotNil(t, after)
	assert.Equal(t, after.ActiveProfile, "foo")

	os.Unsetenv("FABRIC_HOME")
}
