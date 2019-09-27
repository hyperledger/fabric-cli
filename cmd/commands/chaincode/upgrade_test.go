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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

var _ = Describe("ChaincodeUpgradeCommand", func() {
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
		cmd = chaincode.NewChaincodeUpgradeCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode upgrade command", func() {
		Expect(cmd.Name()).To(Equal("upgrade"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("upgrade <chaincode-name> <version> <path>"))
	})
})

var _ = Describe("ChaincodeUpgradeImplementation", func() {
	var (
		impl     *chaincode.UpgradeCommand
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

		impl = &chaincode.UpgradeCommand{}
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

		Context("when chaincode version is not set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
			})

			It("should fail without chaincode version", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode version not specified"))
			})
		})

		Context("when chaincode path is not set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
				impl.ChaincodeVersion = "0.0.0"
			})

			It("should fail without chaincode path", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode path not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
				impl.ChaincodeVersion = "0.0.0"
				impl.ChaincodePath = "path"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChaincodeName = "mycc"
			impl.ChaincodeVersion = "0.0.0"
			impl.ChaincodePath = "github.com/hyperledger/fabric-cli/cmd/commands/chaincode/testdata/chaincode/example"
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		It("should fail without a current context", func() {
			Expect(err).NotTo(BeNil())
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.UpgradeCCReturns(resmgmt.UpgradeCCResponse{}, nil)
			})

			It("should succeed with chaincode upgrade", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal("successfully upgraded chaincode 'mycc'\n"))
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.UpgradeCCReturns(resmgmt.UpgradeCCResponse{}, errors.New("upgrade error"))
			})

			It("should fail to upgrade chaincode", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("upgrade error"))
			})
		})
	})
})
