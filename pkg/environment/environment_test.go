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

var _ = Describe("Profile", func() {
	var (
		settings *environment.Settings
		profile  *environment.Profile
		err      error
	)

	BeforeEach(func() {
		settings = &environment.Settings{}
	})

	Describe("ActiveProfile", func() {
		JustBeforeEach(func() {
			profile, err = settings.GetActiveProfile()
		})

		It("should fail without profiles", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("no profiles currently exist"))
			Expect(profile).To(BeNil())
		})

		Context("when profiles exist", func() {
			BeforeEach(func() {
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}
			})

			It("should fail without active profile set", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("no profile currently active"))
				Expect(profile).To(BeNil())
			})
		})

		Context("when profiles exist but are missing the active profile", func() {
			BeforeEach(func() {
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}
				settings.ActiveProfile = "bar"
			})

			It("should fail to find active profile", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("profile 'bar' was not found"))
				Expect(profile).To(BeNil())
			})
		})

		Context("when active profile exists and is set", func() {
			BeforeEach(func() {
				settings.Profiles = map[string]*environment.Profile{
					"foo": {
						Name: "foo",
					},
				}
				settings.ActiveProfile = "foo"
			})

			It("should return the active profile", func() {
				Expect(err).To(BeNil())
				Expect(profile).NotTo(BeNil())
			})
		})
	})
})
