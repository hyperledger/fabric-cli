/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestMockConfig(t *testing.T) {
	mock := &MockConfig{}
	settings := &environment.Settings{Config: mock}

	temp, err := settings.FromFile()

	assert.NotNil(t, temp)
	assert.Nil(t, err)

	saveErr := temp.Save()

	assert.Nil(t, saveErr)
}

func TestFromFileError(t *testing.T) {
	mock := &MockConfig{}

	mock.ExpectError("config file error")
	settings := &environment.Settings{Config: mock}

	temp, err := settings.FromFile()

	assert.NotNil(t, err)
	assert.Nil(t, temp)
}

func TestSaveError(t *testing.T) {
	mock := &MockConfig{}

	mock.ExpectError("save error")
	settings := &environment.Settings{Config: mock}

	err := settings.Save()

	assert.NotNil(t, err)
}
