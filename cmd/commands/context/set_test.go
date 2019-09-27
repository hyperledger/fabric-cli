/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/hyperledger/fabric-cli/cmd/commands/context"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("SetContextCommand", func() {
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
		cmd = context.NewContextSetCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a set context command", func() {
		Expect(cmd.Name()).To(Equal("set"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("set"))
	})
})

var _ = Describe("SetContextImplementation", func() {
	var (
		impl     *context.SetCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = environment.NewDefaultSettings()
		settings.Home = environment.Home(os.TempDir())
		settings.Streams = environment.Streams{Out: out}

		impl = &context.SetCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail without context details", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when context network is set", func() {
			BeforeEach(func() {
				impl.Context = &environment.Context{
					Network: "foo",
				}
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

		Context("when a context name is provided", func() {
			BeforeEach(func() {
				impl.Name = "bar"
			})

			It("should set the context", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("successfully set context 'bar'"))
			})
		})

		Context("when a context name is not provided", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"bar": {},
					},
					CurrentContext: "bar",
				}
				impl.Context = &environment.Context{
					Organization: "Org1",
				}
			})

			It("should set the context", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("successfully set context 'bar'"))
			})
		})

		Context("when current context is not set", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"bar": {},
					},
				}
			})

			It("should fail to set context", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when context file is invalid", func() {
			BeforeEach(func() {
				impl.Name = "bar"

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

			It("should fail to unmarshal context file", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
