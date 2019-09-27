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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/chaincode"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

var _ = Describe("ChaincodeQueryCommand", func() {
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
		cmd = chaincode.NewChaincodeQueryCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode query command", func() {
		Expect(cmd.Name()).To(Equal("query"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("query <chaincode-name>"))
	})
})

var _ = Describe("ChaincodeQueryImplementation", func() {
	var (
		impl     *chaincode.QueryCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
		factory  *mocks.Factory
		client   *mocks.Channel
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
		client = &mocks.Channel{}

		impl = &chaincode.QueryCommand{}
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

		It("should fail when name is not set", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("chaincode name not specified"))
		})

		Context("when chaincode name is set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
			})

			It("should succeed with chaincode name", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChaincodeName = "mycc"

			impl.Channel = client
		})

		JustBeforeEach(func() {

			err = impl.Run()

		})

		Context("when channel client succeeds", func() {
			BeforeEach(func() {
				client.QueryReturns(channel.Response{}, nil)
			})

			It("should successfully run chaincode query", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when channel client fails", func() {
			BeforeEach(func() {
				client.QueryReturns(channel.Response{}, errors.New("query error"))
			})

			It("should fail to run chaincode query", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query error"))
			})
		})
	})
})
