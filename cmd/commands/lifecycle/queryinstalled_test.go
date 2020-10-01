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
	queryInstalledPlainTextResponse = "Installed chaincodes:\n- Package ID: pkg1, Label: label1\n-- " +
		"References for channel [channel1]:\n--- Name: cc1, Version: v1\n"
	queryInstalledJSONResponse = `[{"package_id":"pkg1","label":"label1","references":{"channel1":[{"name":"cc1","version":"v1"}]}}]`
)

var _ = Describe("LifecycleChaincodeQueryInstalledCommand", func() {
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
		cmd = lifecycle.NewQueryInstalledCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode 'query installed' command", func() {
		Expect(cmd.Name()).To(Equal("queryinstalled"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("queryinstalled"))
	})
})

var _ = Describe("LifecycleChaincodeQueryInstalledImplementation", func() {
	var (
		impl     *lifecycle.QueryInstalledCommand
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

		impl = &lifecycle.QueryInstalledCommand{}
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

		Context("when peer is not set", func() {
			It("should fail without peer", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("peer not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.Peer = "peer1"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.Peer = "peer1"

			result := []resmgmt.LifecycleInstalledCC{
				{
					PackageID: "pkg1",
				},
			}

			client.LifecycleQueryInstalledCCReturns(result, nil)
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.LifecycleQueryInstalledCCReturns([]resmgmt.LifecycleInstalledCC{
					{
						PackageID: "pkg1",
						Label:     "label1",
						References: map[string][]resmgmt.CCReference{
							"channel1": {
								{
									Name:    "cc1",
									Version: "v1",
								},
							},
						},
					},
				}, nil)
			})

			It("should succeed with plain text response", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal(queryInstalledPlainTextResponse))
			})

			When("the output format is set to json", func() {
				BeforeEach(func() {
					impl.OutputFormat = "json"
				})

				It("should succeed with JSON response", func() {
					Expect(err).To(BeNil())
					Expect(fmt.Sprint(out)).To(Equal(queryInstalledJSONResponse))
				})
			})
		})

		Context("when no responses from client", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				client.LifecycleQueryInstalledCCReturns(nil, nil)
			})

			It("should succeed with plain text response", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal("No installed chaincodes on peer peer1"))
			})

			When("the output format is set to json", func() {
				BeforeEach(func() {
					impl.OutputFormat = "json"
				})

				It("should succeed with null JSON response", func() {
					Expect(err).To(BeNil())
					Expect(fmt.Sprint(out)).To(Equal("null"))
				})
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

				client.LifecycleQueryInstalledCCReturns(nil, errors.New("query installed error"))
			})

			It("should fail to query installed chaincodes", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query installed error"))
			})
		})
	})
})
