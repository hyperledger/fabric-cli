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

func TestPluginUinstallCommand(t *testing.T) {
	cmd := NewPluginUninstallCommand(testEnvironment())

	assert.NotNil(t, cmd)
}

func TestUninstallCommandComplete(t *testing.T) {
	pcmd := pluginUninstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
	}

	err := pcmd.complete([]string{"foo"})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.name, "foo")
}

func TestUninstallCommandCompleteError(t *testing.T) {
	pcmd := pluginUninstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
	}

	err := pcmd.complete([]string{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "plugin name not specified")
}

func TestUninstallCommandRun(t *testing.T) {
	pcmd := pluginUninstallCommand{
		out:     new(bytes.Buffer),
		handler: &mocks.MockHandler{},
		name:    "foo",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "successfully uninstalled the plugin\n")
}

func TestUninstallCommandRunErr(t *testing.T) {
	handler := &mocks.MockHandler{}

	handler.ExpectError("an error occurred uninstalling the plugin")

	pcmd := pluginUninstallCommand{
		out:     new(bytes.Buffer),
		handler: handler,
		name:    "foo",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "an error occurred uninstalling the plugin")
}
