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

	pb "github.com/hyperledger/fabric-protos-go/peer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/chaincode"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("ChaincodeListCommand", func() {
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
		cmd = chaincode.NewChaincodeListCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode list command", func() {
		Expect(cmd.Name()).To(Equal("list"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("list"))
	})
})

var _ = Describe("ChaincodeListImplementation", func() {
	var (
		impl     *chaincode.ListCommand
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

		impl = &chaincode.ListCommand{}
		impl.Settings = settings
		impl.Factory = factory
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Run", func() {
		BeforeEach(func() {
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

				client.QueryInstalledChaincodesReturns(&pb.ChaincodeQueryResponse{
					Chaincodes: []*pb.ChaincodeInfo{
						{
							Name: "mycc",
						},
					},
				}, nil)
				client.QueryInstantiatedChaincodesReturns(&pb.ChaincodeQueryResponse{
					Chaincodes: []*pb.ChaincodeInfo{
						{
							Name: "mycc",
						},
					},
				}, nil)
			})

			It("should succeed with chaincode list", func() {
				Expect(err).To(BeNil())

			})
		})

		Context("when resmgmt client fails to list installed chaincode", func() {
			BeforeEach(func() {
				impl.Installed = true

				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.QueryInstalledChaincodesReturns(nil, errors.New("list error"))
			})

			It("should fail to list installed chaincode", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("list error"))
			})
		})

		Context("when resmgmt client fails to list instantiated chaincode", func() {
			BeforeEach(func() {
				impl.Instantiated = true

				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.QueryInstalledChaincodesReturns(nil, errors.New("list error"))
			})

			It("should fail to list instantiated chaincode", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("list error"))
			})
		})
	})
})
