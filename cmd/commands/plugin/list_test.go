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
	plug "github.com/hyperledger/fabric-cli/pkg/plugin"
	"github.com/hyperledger/fabric-cli/pkg/plugin/mocks"
)

var _ = Describe("PluginListCommand", func() {
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
		cmd = plugin.NewPluginListCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a plugin list command", func() {
		Expect(cmd.Name()).To(Equal("list"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
	})
})

var _ = Describe("PluginListImplementation", func() {
	var (
		impl    *plugin.ListCommand
		out     *bytes.Buffer
		handler *mocks.PluginHandler
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		handler = &mocks.PluginHandler{}
	})

	JustBeforeEach(func() {
		impl = &plugin.ListCommand{}
		impl.Settings = &environment.Settings{
			Streams: environment.Streams{
				Out: out,
			},
		}
		impl.Handler = handler
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Run", func() {
		It("should fail if handler fails", func() {
			handler.GetPluginsReturns(nil, errors.New("handler error"))

			err := impl.Run()

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("handler error"))
		})

		It("should succeed when plugins have not been installed", func() {
			Expect(impl.Run()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("no plugins currently exist"))
		})

		Context("when plugins have been installed", func() {
			BeforeEach(func() {
				handler.GetPluginsReturns([]*plug.Plugin{
					{
						Name: "foo",
					},
				}, nil)
			})

			JustBeforeEach(func() {
				err := handler.InstallPlugin("foo")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should successfully list a plugin", func() {
				Expect(impl.Run()).Should(Succeed())
				Expect(fmt.Sprint(out)).To(ContainSubstring("foo\n"))
			})
		})
	})
})
