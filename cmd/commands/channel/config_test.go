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

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("ChannelConfigCommand", func() {
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
		cmd = channel.NewChannelConfigCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a channel config commmand", func() {
		Expect(cmd.Name()).To(Equal("config"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("config"))
	})
})

var _ = Describe("ChannelConfigImplementation", func() {
	var (
		impl     *channel.ConfigCommand
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

		cmd = channel.NewChannelConfigCommand(settings)
		client = &mocks.ResourceManangement{}

		impl = &channel.ConfigCommand{
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
		var (
			cfg *mocks.ChannelCfg
		)

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				cfg = &mocks.ChannelCfg{}
				cfg.IDReturns("mychannel")
				cfg.BlockNumberReturns(0)
				cfg.OrderersReturns([]string{
					"orderer.example.com",
				})
				cfg.AnchorPeersReturns([]*fab.OrgAnchorPeer{
					{
						Host: "peer.example.com",
						Port: 8888,
						Org:  "foo",
					},
				})

				client.QueryConfigFromOrdererReturns(cfg, nil)
			})

			It("should succeed with channel config", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(ContainSubstring("mychannel"))
				Expect(fmt.Sprint(out)).To(ContainSubstring("0"))
				Expect(fmt.Sprint(out)).To(ContainSubstring("orderer.example.com"))
				Expect(fmt.Sprint(out)).To(ContainSubstring("peer.example.com:8888 (foo)"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				client.QueryConfigFromOrdererReturns(nil, errors.New("query error"))
			})

			It("should fail to get channel config", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query error"))
			})
		})
	})
})
