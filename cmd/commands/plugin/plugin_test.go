/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/plugin"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

func TestPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Suite")
}

var _ = Describe("PluginCommand", func() {
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
			cmd = plugin.NewPluginCommand(settings)
		})

		It("should create a plugin command", func() {
			Expect(cmd.Name()).To(Equal("plugin"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("plugin [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("install"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("uninstall"))
		})
	})
})
