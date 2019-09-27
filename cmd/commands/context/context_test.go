/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/context"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Context Suite")
}

var _ = Describe("ContextCommand", func() {
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
			cmd = context.NewContextCommand(settings)
		})

		It("should create a context command", func() {
			Expect(cmd.Name()).To(Equal("context"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("context [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("view"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("use"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("set"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("delete"))
		})
	})
})
