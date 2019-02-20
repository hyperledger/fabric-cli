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

func TestProfileUseCommand(t *testing.T) {
	cmd := NewProfileUseCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.False(t, cmd.HasSubCommands())
}

func TestUseCommandComplete(t *testing.T) {
	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{"foo"})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.name, "foo")
}

func TestUseCommandCompleteError(t *testing.T) {
	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile name not specified")
}

func TestUseCommandCompleteErrorTrim(t *testing.T) {
	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: &environment.Settings{},
	}

	err := pcmd.complete([]string{" "})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile name not specified")
}

func TestUseCommandRun(t *testing.T) {
	mock := &mocks.MockConfig{}
	settings := &environment.Settings{
		Config: mock,
		Profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
	}

	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "successfully set active profile to 'foobar'\n")
}

func TestUseCommandRunError(t *testing.T) {
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

	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "save error")
}

func TestUseCommandRunErrorDNE(t *testing.T) {
	mock := &mocks.MockConfig{}
	settings := &environment.Settings{
		Config: mock,
	}

	pcmd := profileUseCommand{
		out:    new(bytes.Buffer),
		config: settings,
		name:   "foobar",
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "profile 'foobar' was not found")
}
