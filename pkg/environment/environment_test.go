/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var TestPath = filepath.Join(os.TempDir(), "fabric")

func TestEnvironment(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Environment Suite")
}

var _ = Describe("Settings", func() {
	var (
		settings *environment.Settings
	)

	BeforeEach(func() {
		settings = environment.NewDefaultSettings()
	})

	It("should contain default streams", func() {
		Expect(settings.Streams).To(Equal(environment.DefaultStreams))
	})

	Describe("AddFlags", func() {
		var (
			flags *pflag.FlagSet
		)

		BeforeEach(func() {
			flags = pflag.NewFlagSet("test", pflag.ExitOnError)
		})

		JustBeforeEach(func() {
			settings.AddFlags(flags)
		})

		It("should override disable plugins", func() {
			Expect(settings.DisablePlugins).To(BeFalse())

			flags.Set("disable-plugins", "true")

			Expect(settings.DisablePlugins).To(BeTrue())
		})

		It("should override home", func() {
			Expect(settings.Home.String()).To(Equal(environment.DefaultHome.String()))

			flags.Set("home", "path/to/new/home")

			Expect(settings.Home.String()).To(Equal("path/to/new/home"))
		})
	})

	Describe("Init", func() {
		var (
			flags *pflag.FlagSet
			err   error
		)

		BeforeEach(func() {
			flags = pflag.NewFlagSet("test", pflag.ExitOnError)

			settings.AddFlags(flags)
		})

		JustBeforeEach(func() {
			err = settings.Init(flags)
		})

		Context("when environment variables are set", func() {
			BeforeEach(func() {
				os.Setenv("FABRIC_DISABLE_PLUGINS", "true")
			})

			AfterEach(func() {
				os.Unsetenv("FABRIC_DISABLE_PLUGINS")
			})

			It("should disable plugins", func() {
				Expect(err).To(BeNil())
				Expect(settings.DisablePlugins).To(BeTrue())
			})
		})

		Context("when flags are set", func() {
			BeforeEach(func() {
				flags.Set("disable-plugins", "true")
			})

			It("should override disable plugins", func() {
				Expect(err).To(BeNil())
				Expect(settings.DisablePlugins).To(BeTrue())
			})
		})

		Context("when config file exists", func() {
			BeforeEach(func() {
				flags.Set("home", TestPath)

				data, _ := yaml.Marshal(&environment.Config{
					CurrentContext: "baz",
				})

				os.MkdirAll(TestPath, 0777)

				ioutil.WriteFile(filepath.Join(TestPath, environment.DefaultConfigFilename), data, 0777)
			})

			JustAfterEach(func() {
				os.RemoveAll(TestPath)
			})

			It("should populate config", func() {
				Expect(err).To(BeNil())
				Expect(settings.Config.CurrentContext).To(Equal("baz"))
			})
		})

		Context("when config file is invalid", func() {
			BeforeEach(func() {
				flags.Set("home", TestPath)

				data, _ := yaml.Marshal(struct {
					Networks string
				}{
					Networks: "foo",
				})

				os.MkdirAll(TestPath, 0777)

				ioutil.WriteFile(filepath.Join(TestPath, environment.DefaultConfigFilename), data, 0777)
			})

			JustAfterEach(func() {
				os.RemoveAll(TestPath)
			})

			It("should fail to unmarshal config file", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("ModifyConfig", func() {
		var (
			actions []environment.Action
			err     error
		)

		BeforeEach(func() {
			settings.Home = environment.Home(TestPath)
		})

		JustBeforeEach(func() {
			err = settings.ModifyConfig(actions...)
		})

		JustAfterEach(func() {
			os.RemoveAll(TestPath)
		})

		It("should not fail without actions", func() {
			Expect(err).To(BeNil())
		})

		Context("when changing current context", func() {
			BeforeEach(func() {
				actions = append(actions, environment.SetCurrentContext("foo"))
			})

			It("should set current context", func() {
				Expect(err).To(BeNil())
				Expect(settings.Config.CurrentContext).To(Equal("foo"))
			})
		})

		Context("when config file is invalid", func() {
			BeforeEach(func() {
				data, _ := yaml.Marshal(struct {
					Networks string
				}{
					Networks: "foo",
				})

				os.MkdirAll(TestPath, 0777)

				ioutil.WriteFile(filepath.Join(TestPath, environment.DefaultConfigFilename), data, 0777)
			})

			JustAfterEach(func() {
				os.RemoveAll(TestPath)
			})

			It("should fail to unmarshal config file", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("SetupPluginEnvironment", func() {
		JustBeforeEach(func() {
			settings.SetupPluginEnvironment()
		})

		JustAfterEach(func() {
			os.Unsetenv("FABRIC_HOME")
			os.Unsetenv("FABRIC_DISABLE_PLUGINS")
		})

		It("should populate the environment with home", func() {
			v, ok := os.LookupEnv("FABRIC_HOME")

			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(settings.Home.String()))
		})
	})
})
