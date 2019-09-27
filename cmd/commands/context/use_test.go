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

var _ = Describe("UseContextCommand", func() {
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
		cmd = context.NewContextUseCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a set context command", func() {
		Expect(cmd.Name()).To(Equal("use"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("use"))
	})
})

var _ = Describe("SetContextImplementation", func() {
	var (
		impl     *context.UseCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = environment.NewDefaultSettings()
		settings.Home = environment.Home(os.TempDir())
		settings.Streams = environment.Streams{Out: out}

		impl = &context.UseCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail without context name", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when context name is set but does not exist", func() {
			BeforeEach(func() {
				impl.Name = "foo"
			})

			It("should fail without an existing context", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when context name is set", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
				}
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

		Context("when a context exists", func() {
			BeforeEach(func() {
				impl.Name = "bar"
			})

			It("should set the context", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("successfully set current context to 'bar'"))
			})
		})

		Context("when context file is invalid", func() {
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

			It("should fail to unmarshal context file", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
