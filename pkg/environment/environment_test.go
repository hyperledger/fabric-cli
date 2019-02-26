/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

func TestEnvironment(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Environment Suite")
}

var _ = Describe("Environment", func() {
	var (
		settings *environment.Settings
		err      error
	)

	JustBeforeEach(func() {
		settings, err = environment.GetSettings()

		Expect(err).NotTo(HaveOccurred())
	})

	It("should set default home", func() {
		Expect(settings).NotTo(BeNil())
		Expect(settings.Home.String()).To(Equal(environment.DefaultHome.String()))
		Expect(settings.Home.Plugins()).To(Equal(environment.DefaultHome.Plugins()))
	})

	Context("when environment overrides are set", func() {
		BeforeEach(func() {
			err := os.Setenv("FABRIC_HOME", os.TempDir())

			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Unsetenv("FABRIC_HOME")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should override default home", func() {
			Expect(settings).NotTo(BeNil())
			Expect(settings.Home).NotTo(Equal(environment.DefaultHome))
			Expect(settings.Home.String()).To(Equal(os.TempDir()))
			Expect(settings.Home.Plugins()).To(Equal(filepath.Join(os.TempDir(), "plugins")))
		})
	})
})

var _ = Describe("PluginEnvironment", func() {
	var (
		settings *environment.Settings
		err      error
	)

	JustBeforeEach(func() {
		settings, err = environment.GetSettings()

		Expect(err).NotTo(HaveOccurred())

		settings.SetupPluginEnv()
	})

	It("should set environment variables", func() {
		Expect(os.Getenv("FABRIC_HOME")).To(Equal(environment.DefaultHome.String()))
	})

})
