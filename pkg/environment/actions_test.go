/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Actions", func() {
	var (
		config *environment.Config
		action environment.Action
	)

	BeforeEach(func() {
		config = environment.NewConfig()
	})

	JustBeforeEach(func() {
		action(config)
	})

	Describe("SetCurrentContext", func() {
		BeforeEach(func() {
			action = environment.SetCurrentContext("foo")
		})

		It("should set current context", func() {
			Expect(config.CurrentContext).To(Equal("foo"))
		})
	})

	Describe("SetContext", func() {
		Context("when context does not exit", func() {
			BeforeEach(func() {
				action = environment.SetContext("foo", &environment.Context{
					Network:      "foo",
					Organization: "Org1",
					User:         "Admin",
					Channel:      "mychannel",
					Orderers:     []string{"orderer.example.com"},
					Peers:        []string{"peer0.org1.example.com"},
				})
			})

			It("should set context", func() {
				Expect(config.Contexts).To(HaveKey("foo"))
			})
		})

		Context("when context exists", func() {
			BeforeEach(func() {
				config.Contexts["foo"] = &environment.Context{}

				action = environment.SetContext("foo", &environment.Context{
					Network:      "foo",
					Organization: "Org1",
					User:         "Admin",
					Channel:      "mychannel",
					Orderers:     []string{"orderer.example.com"},
					Peers:        []string{"peer0.org1.example.com"},
				})
			})

			It("should set context", func() {
				Expect(config.Contexts).To(HaveKey("foo"))
				Expect(config.Contexts["foo"].Organization).To(Equal("Org1"))
			})
		})
	})

	Describe("DeleteContext", func() {
		BeforeEach(func() {
			config = &environment.Config{
				Contexts: map[string]*environment.Context{
					"foo": {},
				},
			}
			action = environment.DeleteContext("foo")
		})

		It("should set context", func() {
			Expect(config.Contexts).NotTo(HaveKey("foo"))
		})
	})

	Describe("SetNetwork", func() {
		BeforeEach(func() {
			action = environment.SetNetwork("foo", &environment.Network{})
		})

		It("should set network", func() {
			Expect(config.Networks).To(HaveKey("foo"))
		})
	})

	Describe("DeleteNetwork", func() {
		BeforeEach(func() {
			config = &environment.Config{
				Networks: map[string]*environment.Network{
					"foo": {},
				},
			}
			action = environment.DeleteNetwork("foo")
		})

		It("should set context", func() {
			Expect(config.Networks).NotTo(HaveKey("foo"))
		})
	})
})
