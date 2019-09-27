/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/network"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("ViewNetworkCommand", func() {
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
		cmd = network.NewNetworkViewCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a view network command", func() {
		Expect(cmd.Name()).To(Equal("view"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("view"))
	})
})

var _ = Describe("ViewNetworkImplementation", func() {
	var (
		impl     *network.ViewCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = environment.NewDefaultSettings()
		settings.Home = environment.Home(os.TempDir())
		settings.Streams = environment.Streams{Out: out}

		impl = &network.ViewCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err = impl.Run()
		})

		It("should fail without a network", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when a network exists", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					CurrentContext: "foo",
					Contexts: map[string]*environment.Context{
						"foo": {
							Network: "bar",
						},
					},
					Networks: map[string]*environment.Network{
						"bar": {},
					},
				}
			})

			It("should print the current network", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("bar"))
			})
		})

		Context("when a name is provided and found", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Networks: map[string]*environment.Network{
						"foo": {},
					},
				}

				impl.Name = "foo"
			})

			It("should print the specified network", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("foo"))
			})
		})

		Context("when a name is provided but not found", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Networks: map[string]*environment.Network{
						"foo": {},
					},
				}

				impl.Name = "bar"
			})

			It("should print the specified network", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("network 'bar' does not exist"))
			})
		})
	})
})
