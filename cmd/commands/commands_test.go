/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package commands_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Commands Suite")
}

var _ = Describe("Commands", func() {
	var (
		cmds []*cobra.Command
	)

	JustBeforeEach(func() {
		cmds = commands.All(&environment.Settings{})
	})

	It("should not be nil", func() {
		Expect(cmds).NotTo(BeNil())
	})

	It("should contain built in commands", func() {
		Expect(len(cmds)).To(BeNumerically(">", 0))
	})
})
