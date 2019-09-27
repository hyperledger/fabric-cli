/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/chaincode"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

func TestChaincode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chaincode Suite")
}

var _ = Describe("ChaincodeCommand", func() {
	var (
		cmd      *cobra.Command
		settings *environment.Settings
		out      *bytes.Buffer
	)

	Context("when creating a command from settings", func() {
		BeforeEach(func() {
			out = new(bytes.Buffer)

			settings = &environment.Settings{
				Home: environment.Home(os.TempDir()),
				Streams: environment.Streams{
					Out: out,
				},
			}
		})

		JustBeforeEach(func() {
			cmd = chaincode.NewChaincodeCommand(settings)
		})

		It("should create a chaincode command", func() {
			Expect(cmd.Name()).To(Equal("chaincode"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("chaincode [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("install"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("instantiate"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("upgrade"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("query"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("invoke"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("events"))
		})
	})
})

var _ = Describe("BaseChaincodeCommand", func() {
	var c *chaincode.BaseCommand

	BeforeEach(func() {
		c = &chaincode.BaseCommand{}
	})

	Describe("Complete", func() {
		var (
			err           error
			factory       *mocks.Factory
			channelClient *mocks.Channel
			resmgmtClient *mocks.ResourceManagement
		)

		BeforeEach(func() {
			factory = &mocks.Factory{}
			channelClient = &mocks.Channel{}
			resmgmtClient = &mocks.ResourceManagement{}

			factory.ResourceManagementReturns(resmgmtClient, nil)
			factory.ChannelReturns(channelClient, nil)

			c.Factory = factory
		})

		JustBeforeEach(func() {
			err = c.Complete()
		})

		It("should complete", func() {
			Expect(err).To(BeNil())
			Expect(c.Channel).NotTo(BeNil())
		})

		Context("when factory fails to create channel client", func() {
			BeforeEach(func() {
				factory.ChannelReturns(nil, errors.New("factory error"))
			})

			It("should fail with factory error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when factory fails to create resmgmt client", func() {
			BeforeEach(func() {
				factory.ResourceManagementReturns(nil, errors.New("factory error"))
			})

			It("should fail with factory error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
