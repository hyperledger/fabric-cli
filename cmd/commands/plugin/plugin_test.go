/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"bytes"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestPluginCommand(t *testing.T) {
	cmd := NewPluginCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.True(t, cmd.HasSubCommands())
}

func testEnvironment() *environment.Settings {
	return &environment.Settings{
		Home: environment.Home("./tmp"),
		Streams: environment.Streams{
			In:  new(bytes.Buffer),
			Out: new(bytes.Buffer),
			Err: new(bytes.Buffer),
		},
	}
}
