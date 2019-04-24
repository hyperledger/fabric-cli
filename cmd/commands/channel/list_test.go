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

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("ChannelListCommand", func() {
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
		cmd = channel.NewChannelListCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a channel list commmand", func() {
		Expect(cmd.Name()).To(Equal("list"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
	})
})

var _ = Describe("ChannelListImplementation", func() {
	var (
		impl     *channel.ListCommand
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

		cmd = channel.NewChannelListCommand(settings)
		client = &mocks.ResourceManangement{}

		impl = &channel.ListCommand{
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

		Context("when active profile is set", func() {
			BeforeEach(func() {
				settings.ActiveProfile = "foo"
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}
			})

			It("should complete the command", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail without channel id", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("peer not specified"))
		})

		Context("when options are set", func() {
			BeforeEach(func() {
				impl.Options = []resmgmt.RequestOption{
					resmgmt.WithTargets(&peer.Peer{}),
				}
			})

			It("should succeed with channel id set", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				client.QueryChannelsReturns(&pb.ChannelQueryResponse{
					Channels: []*pb.ChannelInfo{
						{
							ChannelId: "mychannel",
						},
					},
				}, nil)
			})

			It("should succeed with channel list", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("mychannel"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				client.QueryChannelsReturns(nil, errors.New("query error"))
			})

			It("should fail to get channel list", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query error"))
			})
		})
	})
})
