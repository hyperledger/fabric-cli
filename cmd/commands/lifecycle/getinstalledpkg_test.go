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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/mocks"
)

var _ = Describe("LifecycleChaincodeGetInstalledPackageCommand", func() {
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
		cmd = lifecycle.NewGetInstalledPkgCommand(settings)
	})

	AfterEach(func() {
		os.Args = args
	})

	It("should create a chaincode getinstalledpackage command", func() {
		Expect(cmd.Name()).To(Equal("getinstalledpackage"))
		Expect(cmd.HasSubCommands()).To(BeFalse())
	})

	It("should provide a help prompt", func() {
		os.Args = append(os.Args, "--help")

		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("getinstalledpackage <peer> <package ID>"))
	})
})

var _ = Describe("LifecycleChaincodeGetInstalledPackageImplementation", func() {
	var (
		impl     *lifecycle.GetInstalledPkgCommand
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

		impl = &lifecycle.GetInstalledPkgCommand{}
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

		It("should fail when peer is not set", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("peer not specified"))
		})

		Context("when package ID is not set", func() {
			BeforeEach(func() {
				impl.Peer = "peer1"
			})

			It("should fail without package ID", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("package ID not specified"))
			})
		})

		Context("when all arguments are set", func() {
			BeforeEach(func() {
				impl.Peer = "peer1"
				impl.PackageID = "pkg1"
			})

			It("should succeed with all arguments", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Run", func() {
		var filePath string

		BeforeEach(func() {
			impl.Peer = "peer1"
			impl.PackageID = "pkg1"
			impl.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
				filePath = filename
				return nil
			}

			result := []byte("pkg1")

			client.LifecycleGetInstalledCCPackageReturns(result, nil)
			impl.ResourceManagement = client
		})

		JustBeforeEach(func() {
			err = impl.Run()
		})

		Context("when failed to write file", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}

				impl.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
					filePath = filename
					return fmt.Errorf("injected write file error")
				}
			})

			It("should fail with writer error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal(fmt.Sprintf("failed to write chaincode package to file %s: injected write file error", filePath)))
			})
		})

		Context("when resmgmt client succeeds", func() {
			BeforeEach(func() {
				settings.Config = &environment.Config{
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
					CurrentContext: "foo",
				}
			})

			It("should succeed with getinstalledpackage", func() {
				Expect(err).To(BeNil())
				Expect(fmt.Sprint(out)).To(Equal(fmt.Sprintf("Chaincode package saved to %s\n", filePath)))
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

				client.LifecycleGetInstalledCCPackageReturns(nil, errors.New("get installed error"))
			})

			It("should fail to get installed chaincode package", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("get installed error"))
			})
		})
	})
})
