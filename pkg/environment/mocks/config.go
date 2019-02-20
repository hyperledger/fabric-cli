/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"errors"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// MockConfig implements the config interface without modifying the filesystem
type MockConfig struct {
	expectErr bool
	errMsg    string
}

// FromFile returns an instance of settings only based on config file
func (m *MockConfig) FromFile() (*environment.Settings, error) {
	if m.expectErr {
		return nil, errors.New(m.errMsg)
	}

	return &environment.Settings{Config: &MockConfig{}}, nil
}

// Save writes the current Settings to the config file
func (m *MockConfig) Save() error {
	if m.expectErr {
		return errors.New(m.errMsg)
	}

	return nil
}

// ExpectError allows tests to mock error
func (m *MockConfig) ExpectError(msg string) {
	m.expectErr = true
	m.errMsg = msg
}
