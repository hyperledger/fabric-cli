/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/network"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

func TestNetwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Network Suite")
}

var _ = Describe("NetworkCommand", func() {
	var (
		cmd      *cobra.Command
		settings *environment.Settings
		out      *bytes.Buffer
	)

	Context("when creating a command from settings", func() {
		BeforeEach(func() {
			out = new(bytes.Buffer)

			settings = &environment.Settings{
				Home: environment.Home(os.TempDir()),
				Streams: environment.Streams{
					Out: out,
				},
			}
		})

		JustBeforeEach(func() {
			cmd = network.NewNetworkCommand(settings)
		})

		It("should create a network command", func() {
			Expect(cmd.Name()).To(Equal("network"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("network [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("view"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("set"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("delete"))
		})
	})
})
