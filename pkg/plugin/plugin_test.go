/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDir = "./testdata/tmp"

func TestGetPlugins(t *testing.T) {
	tests := []struct {
		// test case name
		name string

		// input path
		path string

		// expected output
		expectErr    bool
		outputLength int
	}{
		{
			name:         "Existing Path",
			path:         "./testdata/plugins",
			outputLength: 3,
		},
		{
			name:         "Nonexistent Path",
			path:         "./foo/bar",
			outputLength: 0,
		},
		{
			name:      "Malformed Input",
			path:      "[]",
			expectErr: true,
		},
		{
			name:      "Empty Path",
			path:      "",
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &DefaultHandler{
				Dir:      test.path,
				Filename: DefaultFilename,
			}

			plugins, err := handler.GetPlugins()

			if test.expectErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, plugins)
			assert.Len(t, plugins, test.outputLength)
		})
	}
}

func TestInstallErrors(t *testing.T) {
	defer cleanup()

	tests := []struct {
		name string

		input string

		setup   func()
		cleanup func()
	}{
		{
			name:  "Plugin Does Not Exist",
			input: "./foo/bar",
		},
		{
			name:  "YAML Does Not Exist",
			input: filepath.Join(testDir, "foo", "bar"),
			setup: func() {
				os.MkdirAll(filepath.Join(testDir, "foo", "bar"), 0777)
			},
			cleanup: func() {
				os.RemoveAll(filepath.Join(testDir, "foo"))
			},
		},
		{
			name:  "Malformed YAML",
			input: filepath.Join(testDir, "foo", "bar"),
			setup: func() {
				os.MkdirAll(filepath.Join(testDir, "foo", "bar"), 0777)
				ioutil.WriteFile(filepath.Join(testDir, "foo", "bar", DefaultFilename),
					[]byte("command: !!float 'error'"), 0777,
				)
			},
			cleanup: func() {
				os.RemoveAll(filepath.Join(testDir, "foo"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup()
			}

			if test.cleanup != nil {
				defer test.cleanup()
			}

			handler := &DefaultHandler{
				Dir:      testDir,
				Filename: DefaultFilename,
			}

			err := handler.InstallPlugin(test.input)

			assert.NotNil(t, err)
		})
	}
}

func TestUninstallErrors(t *testing.T) {
	tests := []struct {
		name string

		path   string
		plugin string
	}{
		{
			name:   "Plugin Not Found",
			path:   testDir,
			plugin: "foo",
		},
		{
			name:   "Invalid Path",
			path:   "[]",
			plugin: "foo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &DefaultHandler{
				Dir:      test.path,
				Filename: DefaultFilename,
			}

			err := handler.UninstallPlugin(test.plugin)

			assert.NotNil(t, err)
		})
	}
}

func TestSameFile(t *testing.T) {
	defer cleanup()

	handler := &DefaultHandler{
		Dir:      testDir,
		Filename: DefaultFilename,
	}

	err1 := handler.InstallPlugin("./testdata/plugins/echo")
	assert.Nil(t, err1)

	err2 := handler.UninstallPlugin("echo")
	assert.Nil(t, err2)
}

func TestInstallFileExists(t *testing.T) {
	defer cleanup()

	handler := &DefaultHandler{
		Dir:      testDir,
		Filename: DefaultFilename,
	}

	err1 := handler.InstallPlugin("./testdata/plugins/echo")
	assert.Nil(t, err1)

	err2 := handler.InstallPlugin("./testdata/plugins/echo")
	assert.NotNil(t, err2)
}

func cleanup() {
	os.RemoveAll(testDir)
}
