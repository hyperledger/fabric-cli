/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/hyperledger/fabric-cli/cmd/commands/network"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("NetworkDeleteCommand", func() {
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
		cmd = network.NewNetworkDeleteCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a delete network command", func() {
		Expect(cmd.Name()).To(Equal("delete"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("delete"))
	})
})

var _ = Describe("NetworkDeleteImplementation", func() {
	var (
		impl     *network.DeleteCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = environment.NewDefaultSettings()
		settings.Home = environment.Home(os.TempDir())
		settings.Streams = environment.Streams{Out: out}

		impl = &network.DeleteCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail without network name", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when network name is set", func() {
			BeforeEach(func() {
				impl.Name = "foo"
			})

			It("should successfully validate", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err = impl.Run()
		})

		JustAfterEach(func() {
			os.RemoveAll(impl.Settings.Home.Path(environment.DefaultConfigFilename))
		})

		Context("when a network exists", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Networks: map[string]*environment.Network{
						"foo": {},
						"bar": {},
					},
				}

				impl.Name = "bar"
			})

			It("should delete the network", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("successfully deleted network 'bar'"))
			})
		})

		Context("when network file is invalid", func() {
			BeforeEach(func() {
				data, _ := yaml.Marshal(struct {
					Networks string
				}{
					Networks: "foo",
				})

				os.MkdirAll(settings.Home.String(), 0777)

				ioutil.WriteFile(
					settings.Home.Path(environment.DefaultConfigFilename),
					data,
					0777,
				)
			})

			It("should fail to unmarshal network file", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
