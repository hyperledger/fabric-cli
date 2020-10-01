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
	queryApprovedPlainTextResponse = "Name: cc1, Version: v1, Package ID: pkg1, Sequence: 1, Validation Plugin: vscc," +
		" Endorsement Plugin: escc, Channel Config Policy: ccpolicy, Init Required: true\n- Collection: coll1," +
		" Blocks to Live: 1, Maximum Peer Count: 2, Required Peer Count: 1, MemberOnlyRead: false, cfg.MemberOnlyWrite: false\n"
	queryApprovedJSONResponse = `{"name":"cc1","version":"v1","sequence":1,"endorsement_plugin":"escc","validation_plugin":"vscc"` +
		`,"channel_config_policy":"ccpolicy","collection_config":[{"Payload":{"StaticCollectionConfig":` +
		`{"name":"coll1","required_peer_count":1,"maximum_peer_count":2,"block_to_live":1}}}],"init_required":true,"package_id":"pkg1"}`
)

var _ = Describe("LifecycleChaincodeQueryApprovedCommand", func() {
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
		cmd = lifecycle.NewQueryApprovedCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode 'query approved' command", func() {
		Expect(cmd.Name()).To(Equal("queryapproved"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("queryapproved"))
	})
})

var _ = Describe("LifecycleChaincodeQueryApprovedImplementation", func() {
	var (
		impl     *lifecycle.QueryApprovedCommand
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

		impl = &lifecycle.QueryApprovedCommand{}
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

		Context("when channel is not set", func() {
			It("should fail without channel ID", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("channel ID not specified"))
			})
		})

		Context("when chaincode is not set", func() {
			BeforeEach(func() {
				impl.ChannelID = "channel1"
			})

			It("should fail without chaincode", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode name not specified"))
			})
		})

		Context("when sequence is not set", func() {
			BeforeEach(func() {
				impl.ChannelID = "channel1"
				impl.ChaincodeName = "cc1"
			})

			It("should fail without sequence", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("invalid sequence"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.ChannelID = "channel1"
				impl.ChaincodeName = "cc1"
				impl.Sequence = "1"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChannelID = "channel1"
			impl.ChaincodeName = "cc1"
			impl.Sequence = "1"

			result := resmgmt.LifecycleApprovedChaincodeDefinition{
				Name:                "cc1",
				Version:             "v1",
				Sequence:            1,
				PackageID:           "pkg1",
				ValidationPlugin:    "vscc",
				EndorsementPlugin:   "escc",
				ChannelConfigPolicy: "ccpolicy",
				InitRequired:        true,
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
			}

			client.LifecycleQueryApprovedCCReturns(result, nil)
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
				Expect(fmt.Sprint(out)).To(Equal(queryApprovedPlainTextResponse))
			})

			When("the output format is set to json", func() {
				BeforeEach(func() {
					impl.OutputFormat = "json"
				})

				It("should succeed with JSON response", func() {
					Expect(err).To(BeNil())
					Expect(fmt.Sprint(out)).To(Equal(queryApprovedJSONResponse))
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

				client.LifecycleQueryApprovedCCReturns(resmgmt.LifecycleApprovedChaincodeDefinition{}, errors.New("query approved error"))
			})

			It("should fail to query approved chaincodes", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("query approved error"))
			})
		})
	})
})
