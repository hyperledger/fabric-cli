/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric_test

import (
	"errors"
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
		profile = &environment.Profile{}
		factory = &mocks.Factory{}
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
