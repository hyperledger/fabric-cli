/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/plugin/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPluginInstallCommand(t *testing.T) {
	cmd := NewPluginInstallCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.False(t, cmd.HasSubCommands())
}

func TestInstallCommandComplete(t *testing.T) {
	pcmd := pluginInstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
	}

	err := pcmd.complete([]string{"./foo/bar"})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.path, "./foo/bar")
}

func TestInstallCommandCompleteError(t *testing.T) {
	pcmd := pluginInstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
	}

	err := pcmd.complete([]string{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "plugin path not specified")
}

func TestInstallCommandRun(t *testing.T) {
	pcmd := pluginInstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
		path:    "./foo/bar",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "successfully installed the plugin\n")
}

func TestInstallCommandRunErr(t *testing.T) {
	handler := &mocks.MockHandler{}

	handler.ExpectError("an error occurred installing the plugin")

	pcmd := pluginInstallCommand{
		out:     new(bytes.Buffer),
		handler: handler,
		path:    "./foo/bar",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "an error occurred installing the plugin")
}
