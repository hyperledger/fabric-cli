/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	settings := &environment.Settings{
		Home: environment.DefaultHome,
		Streams: environment.Streams{
			In:  new(bytes.Buffer),
			Out: new(bytes.Buffer),
			Err: new(bytes.Buffer),
		},
	}
	cmd := NewFabricCommand(settings)

	assert.NotNil(t, cmd)
	assert.True(t, cmd.HasSubCommands())

	err := cmd.Execute()

	assert.Nil(t, err)
}

func TestLoadPlugins(t *testing.T) {
	tests := []struct {
		// test name
		name string

		// input arguments
		settings *environment.Settings
		handler  *mocks.MockHandler

		// helper functions
		addPlugin func(handler *mocks.MockHandler, path string) error
		setError  func(handler *mocks.MockHandler, msg string)

		// output
		expectErr bool
		count     int
	}{
		{
			name: "No Plugins",
			settings: &environment.Settings{
				Streams: testStreams(),
			},
			handler: &mocks.MockHandler{},
			count:   0,
		},
		{
			name: "Bad Path",
			settings: &environment.Settings{
				Streams: testStreams(),
			},
			handler: &mocks.MockHandler{},
			setError: func(handler *mocks.MockHandler, msg string) {
				handler.ExpectError(msg)
			},
			expectErr: true,
		},
		{
			name: "One Plugin",
			settings: &environment.Settings{
				Streams: testStreams(),
			},
			handler: &mocks.MockHandler{},
			addPlugin: func(handler *mocks.MockHandler, path string) error {
				return handler.InstallPlugin(path)
			},
			count: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.addPlugin != nil {
				err := test.addPlugin(test.handler, "./foo/bar")
				assert.Nil(t, err)
			}

			if test.setError != nil {
				test.setError(test.handler, "test error message")
			}

			cmd := &cobra.Command{}

			loadPlugins(cmd, test.settings, test.handler)

			err := fmt.Sprint(test.settings.Streams.Err)

			if test.expectErr {
				assert.NotEqual(t, len(err), 0)
				return
			}

			assert.Len(t, err, 0)
			assert.Len(t, cmd.Commands(), test.count)
		})
	}
}

func TestDisablePlugins(t *testing.T) {
	cmd := &cobra.Command{}

	settings := &environment.Settings{
		DisablePlugins: true,
	}

	handler := &mocks.MockHandler{}

	handler.InstallPlugin("./foo/bar")

	loadPlugins(cmd, settings, handler)

	assert.False(t, cmd.HasSubCommands())
}

func testStreams() environment.Streams {
	return environment.Streams{
		In:  new(bytes.Buffer),
		Out: new(bytes.Buffer),
		Err: new(bytes.Buffer),
	}
}
