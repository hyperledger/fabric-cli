/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	cmds := All(&environment.Settings{})

	assert.NotNil(t, cmds)
}
