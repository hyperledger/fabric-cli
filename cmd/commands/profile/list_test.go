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

var _ = Describe("ProfileListCommand", func() {
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
		cmd = profile.NewProfileListCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a profile list commmand", func() {
		Expect(cmd.Name()).To(Equal("list"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
	})
})

var _ = Describe("ProfileListImplementation", func() {
	var (
		impl     *profile.ListCommand
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
			Profiles: make(map[string]*environment.Profile),
		}
	})

	JustBeforeEach(func() {
		impl = &profile.ListCommand{
			Out:      out,
			Settings: settings,
		}
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	It("should fail to list profiles", func() {
		err := impl.Run()

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("no profiles currently exist"))
	})

	Context("when profiles exists", func() {
		JustBeforeEach(func() {
			settings.Profiles = map[string]*environment.Profile{
				"foo": {
					Name: "foo",
				},
				"bar": {
					Name: "bar",
				},
			}

			settings.ActiveProfile = "foo"
		})

		It("should list profiles", func() {
			Expect(impl.Run()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("foo (active)\nbar\n"))
		})
	})
})
