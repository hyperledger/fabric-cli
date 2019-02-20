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

func TestProfileListCommand(t *testing.T) {
	cmd := NewProfileListCommand(testEnvironment())

	assert.NotNil(t, cmd)
	assert.False(t, cmd.HasSubCommands())
}

func TestListCommandRun(t *testing.T) {
	pcmd := &profileListCommand{
		out: new(bytes.Buffer),
		profiles: []*environment.Profile{
			&environment.Profile{
				Name: "foo",
			},
		},
		active: "foo",
	}

	err := pcmd.run()

	assert.Nil(t, err)
	assert.Contains(t, fmt.Sprint(pcmd.out), "foo (active)")
}

func TestListCommandRunError(t *testing.T) {
	pcmd := &profileListCommand{
		out: new(bytes.Buffer),
	}

	err := pcmd.run()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "no profiles currently exist")
}
