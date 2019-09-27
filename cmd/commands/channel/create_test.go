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
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("ChannelCreateCommand", func() {
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
		cmd = channel.NewChannelCreateCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a channel create command", func() {
		Expect(cmd.Name()).To(Equal("create"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("create <channel-id> <tx-path>"))
	})
})

var _ = Describe("ChannelCreateImplementation", func() {
	var (
		impl     *channel.CreateCommand
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

		impl = &channel.CreateCommand{}
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

			It("should fail without channel tx", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("channel tx path not specified"))
			})
		})

		Context("when channel id and tx are set", func() {
			BeforeEach(func() {
				impl.ChannelID = "mychannel"
				impl.ChannelTX = "./testdata/channel.tx"
			})

			It("should succeed with channel id and tx path set", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChannelID = "mychannel"
			impl.ChannelTX = "./testdata/channel.tx"
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				client.SaveChannelReturns(resmgmt.SaveChannelResponse{
					TransactionID: fab.TransactionID("123"),
				}, nil)
			})

			It("should succeed with channel create", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal("successfully created channel 'mychannel'\n"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				client.SaveChannelReturns(resmgmt.SaveChannelResponse{}, errors.New("save error"))
			})

			It("should fail to create channel", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("save error"))
			})
		})
	})
})
