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

	It("should create a channel join command", func() {
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
		factory  *mocks.Factory
		client   *mocks.ResourceManagement
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		factory = &mocks.Factory{}
		client = &mocks.ResourceManagement{}

		impl = &channel.JoinCommand{}
		impl.Settings = settings
		impl.Factory = factory
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
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
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		It("should fail without context", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					CurrentContext: "foo",
					Contexts: map[string]*environment.Context{
						"foo": {
							Peers: []string{"peer0", "peer1"},
						},
					},
				}
				client.JoinChannelReturns(nil)
			})

			It("should succeed with channel join", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal("successfully joined channel 'mychannel'\n"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					CurrentContext: "foo",
					Contexts: map[string]*environment.Context{
						"foo": {
							Peers: []string{"peer0", "peer1"},
						},
					},
				}
				client.JoinChannelReturns(errors.New("join error"))
			})

			It("should fail to get channel join", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("join error"))
			})
		})
	})
})
