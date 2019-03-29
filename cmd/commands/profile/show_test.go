/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/profile"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("ProfileShowCommand", func() {
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
		cmd = profile.NewProfileShowCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a profile show commmand", func() {
		Expect(cmd.Name()).To(Equal("show"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("show [profilename]"))
	})
})

var _ = Describe("ProfileShowImplementation", func() {
	var (
		impl *profile.ShowCommand
		out  *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
	})

	JustBeforeEach(func() {
		impl = &profile.ShowCommand{
			Out: out,
		}
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Complete", func() {
		It("should fail without args", func() {
			err := impl.Complete([]string{})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("no profile currently active"))
		})

		It("should fail with empty string", func() {
			err := impl.Complete([]string{" "})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("profile name not specified"))
		})

		It("should succeed with profile name", func() {
			Expect(impl.Complete([]string{"foo"})).Should(Succeed())
		})

	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err := impl.Complete([]string{"foo"})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail to show profile", func() {
			err := impl.Run()

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("no profiles currently exist"))
		})

		Context("when active profile is set", func() {
			JustBeforeEach(func() {
				impl.Profiles = []*environment.Profile{
					{
						Name: "foo",
					},
				}
				impl.Active = "foo"
			})

			It("should print the active profile", func() {
				Expect(impl.Run()).Should(Succeed())
				Expect(fmt.Sprint(out)).To(ContainSubstring("Name: foo\n"))
			})
		})

		Context("when specified profile does not exist", func() {
			JustBeforeEach(func() {
				impl.Profiles = []*environment.Profile{
					{
						Name: "foo",
					},
				}
				impl.Active = "foo"

				err := impl.Complete([]string{"bar"})

				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail to show profile", func() {
				err := impl.Run()

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("profile 'bar' was not found"))
			})
		})

		Context("when non-active profile is specified", func() {
			JustBeforeEach(func() {
				impl.Profiles = []*environment.Profile{
					{
						Name: "foo",
					},
					{
						Name: "bar",
					},
				}
				impl.Active = "foo"

				err := impl.Complete([]string{"bar"})

				Expect(err).NotTo(HaveOccurred())
			})

			It("should print the active profile", func() {
				Expect(impl.Run()).Should(Succeed())
				Expect(fmt.Sprint(out)).To(ContainSubstring("Name: bar\n"))
			})
		})
	})
})
