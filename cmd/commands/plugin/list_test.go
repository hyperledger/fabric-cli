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

func TestPluginListCommand(t *testing.T) {
	cmd := NewPluginListCommand(testEnvironment())

	assert.NotNil(t, cmd)
}

func TestListCommandRun(t *testing.T) {
	handler := &mocks.MockHandler{}

	handler.InstallPlugin("./foo/bar")

	pcmd := &pluginListCommand{
		out:     new(bytes.Buffer),
		handler: handler,
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Contains(t, fmt.Sprint(pcmd.out), "bar")
}

func TestListCommandRunError(t *testing.T) {
	handler := &mocks.MockHandler{}

	handler.ExpectError("unable to find plugins")

	pcmd := &pluginListCommand{
		out:     new(bytes.Buffer),
		handler: handler,
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "unable to find plugins")
}
