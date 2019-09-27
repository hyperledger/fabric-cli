/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric_test

import (
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/factory.go --fake-name Factory . Factory
//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/channel.go --fake-name Channel . Channel
//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/event.go --fake-name Event . Event
//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/ledger.go --fake-name Ledger . Ledger
//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/resmgmt.go --fake-name ResourceManagement . ResourceManagement
//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/msp.go --fake-name MSP . MSP
//nolint:lll //go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/channelcfg.go --fake-name ChannelCfg github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab.ChannelCfg

func TestFabric(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Fabric Suite")
}

var _ = Describe("Factory", func() {
	var (
		factory fabric.Factory
		err     error

		config *environment.Config
	)

	BeforeEach(func() {
		config = environment.NewConfig()
	})

	JustBeforeEach(func() {
		factory, err = fabric.NewFactory(config)
	})

	It("should fail with empty config", func() {
		Expect(factory).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

	Context("when current context is set", func() {
		BeforeEach(func() {
			config = &environment.Config{
				CurrentContext: "foo",
				Contexts: map[string]*environment.Context{
					"foo": {
						Network: "bar",
					},
				},
			}
		})

		It("should fail to get current network", func() {
			Expect(factory).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})

	Context("when current context and network is set", func() {
		BeforeEach(func() {
			config = &environment.Config{
				CurrentContext: "foo",
				Contexts: map[string]*environment.Context{
					"foo": {
						Network: "bar",
					},
				},
				Networks: map[string]*environment.Network{
					"bar": {},
				},
			}
		})

		It("should create a factory", func() {
			Expect(factory).NotTo(BeNil())
			Expect(err).To(BeNil())
		})
	})

})
