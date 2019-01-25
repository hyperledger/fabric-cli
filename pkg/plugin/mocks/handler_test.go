/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockHandler(t *testing.T) {
	handler := &MockHandler{}

	plugins1, err1 := handler.GetPlugins()

	assert.Nil(t, err1)
	assert.Len(t, plugins1, 0)

	err2 := handler.InstallPlugin("./foo/bar")

	assert.Nil(t, err2)

	plugins2, err3 := handler.GetPlugins()

	assert.Nil(t, err3)
	assert.Len(t, plugins2, 1)

	err4 := handler.UninstallPlugin("bar")

	assert.Nil(t, err4)

	plugins3, err5 := handler.GetPlugins()

	assert.Nil(t, err5)
	assert.Len(t, plugins3, 0)
}

func TestMockHandlerGetError(t *testing.T) {
	handler := &MockHandler{}

	handler.ExpectError("unable to get plugins")

	_, err := handler.GetPlugins()

	assert.NotNil(t, err)
}

func TestMockHandlerInstallError(t *testing.T) {
	handler := &MockHandler{}

	handler.ExpectError("plugin already exists")

	err := handler.InstallPlugin("./foo/bar")

	assert.NotNil(t, err)
}
func TestMockHandlerUninstallError(t *testing.T) {
	handler := &MockHandler{}

	handler.ExpectError("plugin not found")

	err := handler.UninstallPlugin("bar")

	assert.NotNil(t, err)
}
