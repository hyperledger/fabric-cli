/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/profile"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/environment/mocks"
)

var _ = Describe("ProfileDeleteCommand", func() {
	var (
		cmd      *cobra.Command
		settings *environment.Settings
		config   *mocks.DefaultConfig
		out      *bytes.Buffer

		args []string
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		config = &mocks.DefaultConfig{}

		settings = &environment.Settings{
			Config: config,
			Home:   environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		config.FromFileReturns(settings, nil)
		config.SaveReturns(nil)

		args = os.Args
	})

	JustBeforeEach(func() {
		cmd = profile.NewProfileDeleteCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a profile delete commmand", func() {
		Expect(cmd.Name()).To(Equal("delete"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("delete <profilename>"))
	})
})

var _ = Describe("ProfileDeleteImplementation", func() {
	var (
		impl     *profile.DeleteCommand
		out      *bytes.Buffer
		settings *environment.Settings
		config   *mocks.DefaultConfig
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		config = &mocks.DefaultConfig{}

		settings = &environment.Settings{
			Config: config,
			Home:   environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		config.FromFileReturns(settings, nil)
		config.SaveReturns(nil)
	})

	JustBeforeEach(func() {
		impl = &profile.DeleteCommand{
			Out:      out,
			Settings: settings,
		}
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Complete", func() {
		It("should fail without args", func() {
			err := impl.Complete([]string{})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("profile name not specified"))
		})

		It("should fail with empty string", func() {
			err := impl.Complete([]string{" "})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("profile name not specified"))
		})

		It("should succeed with profile name", func() {
			Expect(impl.Complete([]string{"foo"})).Should(Succeed())
		})

		Context("when config cannot be loaded", func() {
			BeforeEach(func() {
				config.FromFileReturns(nil, errors.New("cannot load config"))
			})

			It("should fail loading config", func() {
				err := impl.Complete([]string{"foo"})

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("cannot load config"))
			})
		})
	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err := impl.Complete([]string{"foo"})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail with non-existent profile", func() {
			err := impl.Run()

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("profile 'foo' was not found"))
		})

		Context("when a profile exists", func() {
			BeforeEach(func() {
				settings.Profiles = []*environment.Profile{
					{
						Name: "foo",
					},
				}
				settings.ActiveProfile = "foo"
			})

			It("should successfully delete the profile", func() {
				Expect(impl.Run()).Should(Succeed())
				Expect(fmt.Sprint(out)).To(ContainSubstring("successfully deleted profile 'foo'\n"))
				Expect(settings.ActiveProfile).To(Equal(""))
			})
		})

		Context("when config cannot be saved", func() {
			BeforeEach(func() {
				config.SaveReturns(errors.New("save error"))

				settings.Profiles = []*environment.Profile{
					{
						Name: "foo",
					},
				}
			})

			It("should fail loading config", func() {
				err := impl.Run()

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("save error"))
			})
		})
	})
})
