/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("LifecycleChaincodePackageCommand", func() {
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
		cmd = lifecycle.NewPackageCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode package command", func() {
		Expect(cmd.Name()).To(Equal("package"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("package <chaincode-label> <chaincode-type> <path>"))
	})
})

var _ = Describe("LifecycleChaincodePackageImplementation", func() {
	var (
		impl     *lifecycle.PackageCommand
		err      error
		out      *bytes.Buffer
		settings *environment.Settings
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		impl = &lifecycle.PackageCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail when label is not set", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("chaincode label not specified"))
		})

		Context("when chaincode path is not set", func() {
			BeforeEach(func() {
				impl.Label = "mycc"
			})

			It("should fail without chaincode path", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode path not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.Label = "mycc"
				impl.Path = "path"
				impl.Type = "golang"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.Label = "mycc"
			impl.Path = "github.com/hyperledger/fabric-cli/cmd/commands/lifecycle/testdata/chaincode/example/example.go"
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when chaincode path is invalid", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				impl.Path = "path/to/chaincode"
			})

			It("should fail with an invalid path", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
