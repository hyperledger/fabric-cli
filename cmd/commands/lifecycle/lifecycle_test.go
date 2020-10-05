/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

func TestLifecycle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifecycle Suite")
}

var _ = Describe("LifecycleChaincodeCommand", func() {
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
			cmd = lifecycle.NewCommand(settings)
		})

		It("should create a lifecycle command", func() {
			Expect(cmd.Name()).To(Equal("lifecycle"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("lifecycle [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("package"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("install"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("approve"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("commit"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("queryinstalled"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("queryapproved"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("checkcommitreadiness"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("querycommitted"))
		})
	})
})

var _ = Describe("LifecycleBaseChaincodeCommand", func() {
	var c *lifecycle.BaseCommand

	BeforeEach(func() {
		c = &lifecycle.BaseCommand{}
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
