/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package profile

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/environment/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProfileDeleteCommand(t *testing.T) {
	cmd := NewProfileDeleteCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.False(t, cmd.HasSubCommands())
}

func TestDeleteCommandComplete(t *testing.T) {
	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{"foobar"})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.name, "foobar")
}

func TestDeleteCommandCompleteError(t *testing.T) {
	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile name not specified")
}

func TestDeleteCommandCompleteErrorTrim(t *testing.T) {
	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{" "})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile name not specified")
}

func TestDeleteCommandRun(t *testing.T) {
	mock := &mocks.MockConfig{}
	settings := &environment.Settings{
		Config: mock,
		Profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
		ActiveProfile: "foobar",
	}

	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "successfully deleted profile 'foobar'\n")
	assert.Len(t, settings.ActiveProfile, 0)
}

func TestDeleteCommandRunError(t *testing.T) {
	mock := &mocks.MockConfig{}
	settings := &environment.Settings{
		Config: mock,
		Profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
	}

	mock.ExpectError("save error")

	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "save error")
}

func TestDeleteCommandErrorDNE(t *testing.T) {
	mock := &mocks.MockConfig{}
	settings := &environment.Settings{
		Config: mock,
	}

	pcmd := profileDeleteCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile 'foobar' was not found")
}
