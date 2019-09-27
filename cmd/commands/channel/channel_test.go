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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

func TestChannel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Channel Suite")
}

var _ = Describe("ChannelCommand", func() {
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
			cmd = channel.NewChannelCommand(settings)
		})

		It("should create a channel command", func() {
			Expect(cmd.Name()).To(Equal("channel"))
			Expect(cmd.HasSubCommands()).To(BeTrue())
			Expect(cmd.Execute()).Should(Succeed())
			Expect(fmt.Sprint(out)).To(ContainSubstring("channel [command]"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("create"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("join"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("update"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
			Expect(fmt.Sprint(out)).To(ContainSubstring("config"))
		})
	})
})
var _ = Describe("BaseChannelCommand", func() {
	var c *channel.BaseCommand

	BeforeEach(func() {
		c = &channel.BaseCommand{}
	})

	Describe("Complete", func() {
		var (
			err     error
			factory *mocks.Factory
			client  *mocks.ResourceManagement
		)

		BeforeEach(func() {
			factory = &mocks.Factory{}
			client = &mocks.ResourceManagement{}

			factory.ResourceManagementReturns(client, nil)

			c.Factory = factory
		})

		JustBeforeEach(func() {
			err = c.Complete()
		})

		It("should complete", func() {
			Expect(err).To(BeNil())
			Expect(c.ResourceManagement).NotTo(BeNil())
		})

		Context("when resmgmt fails", func() {
			BeforeEach(func() {
				factory.ResourceManagementReturns(nil, errors.New("factory error"))
			})

			It("should fail with factory error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
