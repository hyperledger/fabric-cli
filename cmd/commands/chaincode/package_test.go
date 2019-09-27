/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/chaincode"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

var _ = Describe("ChaincodePackageCommand", func() {
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
		cmd = chaincode.NewChaincodePackageCommand(settings)
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
		Expect(fmt.Sprint(out)).To(ContainSubstring("package <chaincode-name> <path>"))
	})
})

var _ = Describe("ChaincodePackageImplementation", func() {
	var (
		impl     *chaincode.PackageCommand
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

		impl = &chaincode.PackageCommand{}
		impl.Settings = settings
	})

	It("should not be nil", func() {
		Expect(impl).ShouldNot(BeNil())
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = impl.Validate()
		})

		It("should fail when name is not set", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("chaincode name not specified"))
		})

		Context("when chaincode path is not set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
			})

			It("should fail without chaincode path", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("chaincode path not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.ChaincodeName = "mycc"
				impl.ChaincodePath = "path"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			impl.ChaincodeName = "mycc"
			impl.ChaincodePath = "github.com/hyperledger/fabric-cli/cmd/commands/chaincode/testdata/chaincode/example/example.go"
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

				impl.ChaincodePath = "path/to/chaincode"
			})

			It("should fail with an invalid path", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
