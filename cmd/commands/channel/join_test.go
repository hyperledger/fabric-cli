/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("ChannelJoinCommand", func() {
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
		cmd = channel.NewChannelJoinCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a channel join commmand", func() {
		Expect(cmd.Name()).To(Equal("join"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("join <channel-id>"))
	})
})

var _ = Describe("ChannelJoinImplementation", func() {
	var (
		impl     *channel.JoinCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
		cmd      *cobra.Command
		client   *mocks.ResourceManangement
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		cmd = channel.NewChannelJoinCommand(settings)
		client = &mocks.ResourceManangement{}

		impl = &channel.JoinCommand{
			Out:                 out,
			Settings:            settings,
			ResourceManangement: client,
		}
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Complete", func() {
		JustBeforeEach(func() {
			err = impl.Complete(cmd)
		})

		It("should fail without profiles", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("no profiles currently exist"))
		})

		Context("when args are provided", func() {
			BeforeEach(func() {
				settings.ActiveProfile = "foo"
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}

				cmd.Flags().Parse([]string{"mychannel"})
			})

			It("should populate channel id", func() {
				Expect(err).To(BeNil())
				Expect(impl.ChannelID).To(Equal("mychannel"))
			})
		})

		Context("when too many args are provided", func() {
			BeforeEach(func() {
				settings.ActiveProfile = "foo"
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}

				cmd.Flags().Parse([]string{"foo", "bar", "baz"})
			})

			It("should fail to complete", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("unexpected args"))
			})
		})
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail without channel id", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("channel id not specified"))
		})

		Context("when channel id is set", func() {
			BeforeEach(func() {
				impl.ChannelID = "mychannel"
			})

			It("should succeed with channel id set", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChannelID = "mychannel"
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				client.JoinChannelReturns(nil)
			})

			It("should succeed with channel join", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal("successfully joined channel 'mychannel'\n"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				client.JoinChannelReturns(errors.New("join error"))
			})

			It("should fail to get channel join", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("join error"))
			})
		})
	})
})
