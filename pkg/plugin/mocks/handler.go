/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"errors"
	"path/filepath"

	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// MockHandler is a mock plugin handler
type MockHandler struct {
	plugins []*plugin.Plugin

	expectErr bool
	errMsg    string
}

// GetPlugins returns all installed plugins found in the plugins directory
func (m *MockHandler) GetPlugins() ([]*plugin.Plugin, error) {
	if m.expectErr {
		return nil, errors.New(m.errMsg)
	}

	return m.plugins, nil
}

// InstallPlugin creates a symlink from the specified directory to the plugins
// directory
func (m *MockHandler) InstallPlugin(path string) error {
	if m.expectErr {
		return errors.New(m.errMsg)
	}

	m.plugins = append(m.plugins, &plugin.Plugin{
		Name: filepath.Base(path),
	})

	return nil
}

// UninstallPlugin removes the symlink from the plugin directory by plugin name
func (m *MockHandler) UninstallPlugin(name string) error {
	if m.expectErr {
		return errors.New(m.errMsg)
	}

	for i, plugin := range m.plugins {
		if plugin.Name == name {
			m.plugins = append(m.plugins[:i], m.plugins[i+1:]...)
		}
	}

	return nil
}

// ExpectError allows tests to mock error
func (m *MockHandler) ExpectError(msg string) {
	m.expectErr = true
	m.errMsg = msg
}
