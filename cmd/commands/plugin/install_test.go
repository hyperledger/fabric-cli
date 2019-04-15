/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/plugin"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin/mocks"
)

var _ = Describe("PluginInstallCommand", func() {
	var (
		cmd      *cobra.Command
		settings *environment.Settings
		out      *bytes.Buffer

		args []string
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		args = os.Args
	})

	JustBeforeEach(func() {
		cmd = plugin.NewPluginInstallCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a plugin install commmand", func() {
		Expect(cmd.Name()).To(Equal("install"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("install <plugin-path>"))
	})
})

var _ = Describe("PluginInstallImplementation", func() {
	var (
		impl    *plugin.InstallCommand
		out     *bytes.Buffer
		handler *mocks.PluginHandler
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		handler = &mocks.PluginHandler{}
	})

	JustBeforeEach(func() {
		impl = &plugin.InstallCommand{
			Out:     out,
			Handler: handler,
		}
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Complete", func() {
		It("should fail without args", func() {
			err := impl.Complete([]string{})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("plugin path not specified"))
		})

		It("should fail with empty string", func() {
			err := impl.Complete([]string{" "})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("plugin path not specified"))
		})

		It("should succeed with input path", func() {
			Expect(impl.Complete([]string{"./foo"})).Should(Succeed())
		})
	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err := impl.Complete([]string{"./foo"})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail if handler fails", func() {
			handler.InstallPluginReturns(errors.New("handler error"))

			err := impl.Run()

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("handler error"))
		})

		It("should successfully install a plugin", func() {
			Expect(impl.Run()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("successfully installed the plugin\n"))
		})
	})
})
