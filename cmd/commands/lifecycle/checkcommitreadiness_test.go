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

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

const (
	checkCommitReadinessPlainTextResponse = "Approving orgs: [org1]\nNon-approving orgs: [org2]\n"
	checkCommitReadinessJSONResponse      = `{"approvals":{"org1":true,"org2":false}}`
)

var _ = Describe("LifecycleChaincodeCheckCommitReadinessCommand", func() {
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
		cmd = lifecycle.NewCheckCommitReadinessCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode checkcommitreadiness command", func() {
		Expect(cmd.Name()).To(Equal("checkcommitreadiness"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("checkcommitreadiness"))
	})
})

var _ = Describe("LifecycleChaincodeCheckCommitReadinessImplementation", func() {
	var (
		impl     *lifecycle.CheckCommitReadinessCommand
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

		impl = &lifecycle.CheckCommitReadinessCommand{}
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

		It("should fail when chaincode name is not set", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("chaincode name not specified"))
		})

		Context("when chaincode version is not set", func() {
			BeforeEach(func() {
				impl.Name = "mycc"
			})

			It("should fail without chaincode version", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode version not specified"))
			})
		})

		Context("when chaincode sequence is not set", func() {
			BeforeEach(func() {
				impl.Name = "mycc"
				impl.Version = "0.0.0"
			})

			It("should fail without chaincode sequence", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("sequence not specified"))
			})
		})

		Context("when chaincode sequence is not greater than 0", func() {
			BeforeEach(func() {
				impl.Name = "mycc"
				impl.Version = "0.0.0"
				impl.Sequence = "-1"
			})

			It("should fail with chaincode sequence not greater than 0", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("sequence must be greater than 0"))
			})
		})

		Context("when chaincode sequence is invalid", func() {
			BeforeEach(func() {
				impl.Name = "mycc"
				impl.Version = "0.0.0"
				impl.Sequence = "xxx"
			})

			It("should fail with chaincode sequence is invalid", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("invalid sequence"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.Name = "mycc"
				impl.Version = "0.0.0"
				impl.Sequence = "1"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.Name = "cc1"
			impl.Sequence = "1"

			result := resmgmt.LifecycleCheckCCCommitReadinessResponse{
				Approvals: map[string]bool{"org1": true, "org2": false},
			}

			client.LifecycleCheckCCCommitReadinessReturns(result, nil)
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {
							Peers: []string{"peer1"},
						},
					},
					CurrentContext: "foo",
				}
			})

			It("should succeed with plain text response", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal(checkCommitReadinessPlainTextResponse))
			})

			When("the output format is set to json", func() {
				BeforeEach(func() {
					impl.OutputFormat = "json"
				})

				It("should succeed with JSON response", func() {
					Expect(err).To(BeNil())
					Expect(fmt.Sprint(out)).To(Equal(checkCommitReadinessJSONResponse))
				})
			})
		})

		Context("when resmgmt client fails", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {
							Peers: []string{"peer1"},
						},
					},
					CurrentContext: "foo",
				}

				client.LifecycleCheckCCCommitReadinessReturns(resmgmt.LifecycleCheckCCCommitReadinessResponse{}, errors.New("check commit readiness error"))
			})

			It("should fail to check chaincode for commit readiness", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("check commit readiness error"))
			})
		})
	})
})
