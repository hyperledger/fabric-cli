/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

//go:generate counterfeiter -o mocks/factory.go --fake-name Factory . Factory
//go:generate counterfeiter -o mocks/channel.go --fake-name Channel . Channel
//go:generate counterfeiter -o mocks/event.go --fake-name Event . Event
//go:generate counterfeiter -o mocks/ledger.go --fake-name Ledger . Ledger
//go:generate counterfeiter -o mocks/resmgmt.go --fake-name ResourceManangement . ResourceManagement
//go:generate counterfeiter -o mocks/msp.go --fake-name MSP . MSP

func TestFabric(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Fabric Suite")
}

var _ = Describe("Fabric", func() {
	var (
		profile *environment.Profile
		factory *mocks.Factory
	)

	BeforeEach(func() {
		profile = &environment.Profile{
			Context: &environment.Context{
				Organization: "Org1",
				Identity:     "Admin",
				Orderers:     []string{"orderer.example.com"},
				Peers:        []string{"peer0.org1.example.com"},
			},
			CryptoConfig:    "${FABRIC_CFG_PATH}",
			CredentialStore: filepath.Clean(os.TempDir()),
			Channels: map[string]*environment.Channel{
				"mychannel": {
					ID:    "mychannel",
					Peers: []string{"peer0.org1.example.com"},
				},
			},
			Organizations: map[string]*environment.Organization{
				"Org1": {
					ID: "Org1",
					MSP: &environment.MSP{
						ID:    "Org1MSP",
						Store: "peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp",
					},
					Peers: []string{
						"peer0.example.com",
					},
				},
			},
			Orderers: map[string]*environment.Orderer{
				"orderer.example.com": {
					ID:  "orderer.example.com",
					URL: "grpc://localhost:7050",
					TLS: "${FABRIC_CFG_PATH}/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
				},
			},
			Peers: map[string]*environment.Peer{
				"peer0.org1.example.com": {
					ID:  "orderer.example.com",
					URL: "grpc://localhost:7050",
					TLS: "${FABRIC_CFG_PATH}/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem",
					ChannelOptions: map[string]interface{}{
						"endorsingPeer": true,
					},
					GRPCOptions: map[string]interface{}{
						"keep-alive-timeout": "20s",
					},
				},
			},
		}

		factory = &mocks.Factory{}
	})

	It("should transform profile to config", func() {
		config, err := profile.ToTemplate("./templates/config.tmpl")

		Expect(err).To(BeNil())
		Expect(config).NotTo(BeNil())
	})

	Describe("Channel", func() {
		var (
			client *fabric.ChannelClient
			err    error
		)

		JustBeforeEach(func() {
			client, err = fabric.NewChannelClient(profile, fabric.WithFactory(factory))
		})

		It("should not be nil", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})

		Context("when factory fails to create client", func() {
			BeforeEach(func() {
				factory.ChannelReturns(nil, errors.New("initialization error"))
			})

			It("should fail to create client", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("initialization error"))
			})
		})
	})

	Describe("Event", func() {
		var (
			client *fabric.EventClient
			err    error
		)

		JustBeforeEach(func() {
			client, err = fabric.NewEventClient(profile, fabric.WithFactory(factory))
		})

		It("should not be nil", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})

		Context("when factory fails to create client", func() {
			BeforeEach(func() {
				factory.EventReturns(nil, errors.New("initialization error"))
			})

			It("should fail to create client", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("initialization error"))
			})
		})
	})

	Describe("Ledger", func() {
		var (
			client *fabric.LedgerClient
			err    error
		)

		JustBeforeEach(func() {
			client, err = fabric.NewLedgerClient(profile, fabric.WithFactory(factory))
		})

		It("should not be nil", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})

		Context("when factory fails to create client", func() {
			BeforeEach(func() {
				factory.LedgerReturns(nil, errors.New("initialization error"))
			})

			It("should fail to create client", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("initialization error"))
			})
		})
	})

	Describe("ResourceManagement", func() {
		var (
			client *fabric.ResourceManagementClient
			err    error
		)

		JustBeforeEach(func() {
			client, err = fabric.NewResourceManagementClient(profile, fabric.WithFactory(factory))
		})

		It("should not be nil", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})

		Context("when factory fails to create client", func() {
			BeforeEach(func() {
				factory.ResourceManagementReturns(nil, errors.New("initialization error"))
			})

			It("should fail to create client", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("initialization error"))
			})
		})
	})

	Describe("MSP", func() {
		var (
			client *fabric.MSPClient
			err    error
		)

		JustBeforeEach(func() {
			client, err = fabric.NewMSPClient(profile, fabric.WithFactory(factory))
		})

		It("should not be nil", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
		})
	})
})
