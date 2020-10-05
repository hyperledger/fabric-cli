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

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

const (
	queryCommittedPlainTextResponse = "Name: cc1, Version: v1, Sequence: 1, Validation Plugin: vscc," +
		" Endorsement Plugin: escc, Channel Config Policy: ccpolicy, Init Required: true," +
		" Approving orgs: [org1], Non-approving orgs: [org2]\n- Collection: coll1," +
		" Blocks to Live: 1, Maximum Peer Count: 2, Required Peer Count: 1, MemberOnlyRead: false, cfg.MemberOnlyWrite: false\n"
	queryCommittedJSONResponse = `[{"name":"cc1","version":"v1","sequence":1,"endorsement_plugin":"escc","validation_plugin":"vscc"` +
		`,"channel_config_policy":"ccpolicy","collection_config":[{"Payload":{"StaticCollectionConfig":{"name":"coll1",` +
		`"required_peer_count":1,"maximum_peer_count":2,"block_to_live":1}}}],"init_required":true,"approvals":{"org1":true,"org2":false}}]`
)

var _ = Describe("LifecycleChaincodeQueryCommittedCommand", func() {
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
		cmd = lifecycle.NewQueryCommittedCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode 'query committed' command", func() {
		Expect(cmd.Name()).To(Equal("querycommitted"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("querycommitted"))
	})
})

var _ = Describe("LifecycleChaincodeQueryCommittedImplementation", func() {
	var (
		impl     *lifecycle.QueryCommittedCommand
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

		impl = &lifecycle.QueryCommittedCommand{}
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

		Context("when chaincode is not set", func() {
			It("should fail without chaincode", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode name not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "cc1"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChaincodeName = "cc1"

			result := resmgmt.LifecycleChaincodeDefinition{
				Name:                "cc1",
				Version:             "v1",
				Sequence:            1,
				EndorsementPlugin:   "escc",
				ValidationPlugin:    "vscc",
				ChannelConfigPolicy: "ccpolicy",
				CollectionConfig: []*pb.CollectionConfig{
					{
						Payload: &pb.CollectionConfig_StaticCollectionConfig{
							StaticCollectionConfig: &pb.StaticCollectionConfig{
								Name:              "coll1",
								RequiredPeerCount: 1,
								MaximumPeerCount:  2,
								BlockToLive:       1,
							},
						},
					},
				},
				InitRequired: true,
				Approvals:    map[string]bool{"org1": true, "org2": false},
			}

			client.LifecycleQueryCommittedCCReturns([]resmgmt.LifecycleChaincodeDefinition{result}, nil)
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
				Expect(fmt.Sprint(out)).To(Equal(queryCommittedPlainTextResponse))
			})

			When("the output format is set to json", func() {
				BeforeEach(func() {
					impl.OutputFormat = "json"
				})

				It("should succeed with JSON response", func() {
					Expect(err).To(BeNil())
					Expect(fmt.Sprint(out)).To(Equal(queryCommittedJSONResponse))
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

				client.LifecycleQueryCommittedCCReturns(nil, errors.New("query committed error"))
			})

			It("should fail to query committed chaincodes", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query committed error"))
			})
		})
	})
})
