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
	"github.com/stretchr/testify/assert"
)

func TestProfileShowCommand(t *testing.T) {
	cmd := NewProfileShowCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.False(t, cmd.HasSubCommands())
}
func TestShowCommandCompleteDefault(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
		active: "foobar",
	}

	err := pcmd.complete([]string{})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.name, "foobar")
}

func TestShowCommandCompleteDifferent(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
			&environment.Profile{
				Name: "baz",
			},
		},
		active: "foobar",
	}

	err := pcmd.complete([]string{"baz"})

	assert.Nil(t, err)
	assert.Equal(t, pcmd.name, "baz")
}

func TestShowCommandCompleteError(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
	}

	err := pcmd.complete([]string{})

	assert.NotNil(t, err)
}

func TestShowCommandRunDefault(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
		},
		active: "foobar",
		name:   "foobar",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "Name: foobar\n")
}

func TestShowCommandRunDifferent(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foobar",
			},
			&environment.Profile{
				Name: "baz",
			},
		},
		active: "foobar",
		name:   "baz",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(pcmd.out), "Name: baz\n")
}

func TestShowCommandRunError(t *testing.T) {
	pcmd := profileShowCommand{
		out: new(bytes.Buffer),
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "no profiles currently exist")
}
