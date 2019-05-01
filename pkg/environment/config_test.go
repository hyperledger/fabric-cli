/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Context", func() {
	var context *environment.Context

	BeforeEach(func() {
		context = &environment.Context{}
	})

	Describe("ToString", func() {
		var contextString string

		JustBeforeEach(func() {
			contextString = context.String()
		})

		It("should return a string", func() {
			Expect(contextString).NotTo(BeEmpty())
			Expect(contextString).To(ContainSubstring("Network"))
		})
	})
})

var _ = Describe("Network", func() {
	var network *environment.Network

	BeforeEach(func() {
		network = &environment.Network{}
	})

	Describe("ToString", func() {
		var networkString string

		JustBeforeEach(func() {
			networkString = network.String()
		})

		It("should return a string", func() {
			Expect(networkString).NotTo(BeEmpty())
			Expect(networkString).To(ContainSubstring("Path"))
		})
	})
})

var _ = Describe("Config", func() {
	var (
		config *environment.Config
	)

	BeforeEach(func() {
		config = environment.NewConfig()
	})

	It("should contain default settings", func() {
		Expect(config).NotTo(BeNil())
	})

	Describe("AddFlags", func() {
		var (
			flags *pflag.FlagSet
		)

		BeforeEach(func() {
			flags = pflag.NewFlagSet("test", pflag.ExitOnError)
		})

		JustBeforeEach(func() {
			config.AddFlags(flags)
		})

		Context("when current context set", func() {
			BeforeEach(func() {
				config.CurrentContext = "foo"
			})

			It("should not erase current context", func() {
				Expect(config.CurrentContext).To(Equal("foo"))
			})

			It("should override current context with flag", func() {
				flags.Set("context", "bar")

				Expect(config.CurrentContext).To(Equal("bar"))
			})
		})
	})

	Describe("LoadFromFile", func() {
		var (
			config *environment.Config
			err    error

			path = TestPath
		)

		BeforeEach(func() {
			config = environment.NewConfig()

			os.MkdirAll(path, 0777)
		})

		JustBeforeEach(func() {
			err = config.LoadFromFile(filepath.Join(path, environment.DefaultConfigFilename))
		})

		JustAfterEach(func() {
			os.RemoveAll(path)
		})

		It("should fail to find config file", func() {
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		Context("when config file exists", func() {
			BeforeEach(func() {
				data, _ := yaml.Marshal(&environment.Config{
					Networks: map[string]*environment.Network{
						"foo.bar": {
							ConfigPath: "path/to/sdk",
						},
					},
					Contexts: map[string]*environment.Context{
						"baz": {
							User: "Admin",
						},
					},
					CurrentContext: "baz",
				})

				ioutil.WriteFile(filepath.Join(path, environment.DefaultConfigFilename), data, 0777)
			})

			It("should populate current context", func() {
				Expect(config.CurrentContext).To(Equal("baz"))
			})

			It("should handle map keys with dots", func() {
				Expect(config.Networks).To(HaveKey("foo.bar"))
			})

			It("should populate context baz", func() {
				Expect(config.Contexts).To(HaveKey("baz"))
				Expect(config.Contexts["baz"].User).To(Equal("Admin"))
			})
		})

		Context("when config file is invalid", func() {
			BeforeEach(func() {
				data, _ := yaml.Marshal(struct {
					Networks string
				}{
					Networks: "foo",
				})

				ioutil.WriteFile(filepath.Join(path, environment.DefaultConfigFilename), data, 0777)
			})

			It("should fail to unmarshal config file", func() {
				Expect(err).NotTo(BeNil())
			})

		})
	})

	Describe("Save", func() {
		var (
			config *environment.Config
			err    error

			path = TestPath
		)

		BeforeEach(func() {
			config = environment.NewConfig()

			os.MkdirAll(path, 0777)
		})

		JustBeforeEach(func() {
			err = config.Save(filepath.Join(path, environment.DefaultConfigFilename))
		})

		JustAfterEach(func() {
			os.RemoveAll(path)
		})

		Context("when config is populated", func() {
			BeforeEach(func() {
				config = &environment.Config{
					Networks: map[string]*environment.Network{
						"foo.bar": {
							ConfigPath: "path/to/sdk",
						},
					},
					Contexts: map[string]*environment.Context{
						"baz": {
							User: "Admin",
						},
					},
					CurrentContext: "baz",
				}
			})

			It("should write the config file", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetCurrentContext", func() {
		var (
			context *environment.Context
			err     error
		)
		JustBeforeEach(func() {
			context, err = config.GetCurrentContext()
		})

		It("should fail without a current context", func() {
			Expect(context).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		Context("when current context doesnt not exist", func() {
			BeforeEach(func() {
				config = &environment.Config{
					CurrentContext: "foo",
				}
			})

			It("should fail to return current context", func() {
				Expect(context).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when current context is set and exits", func() {
			BeforeEach(func() {
				config = &environment.Config{
					CurrentContext: "foo",
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
				}
			})

			It("should return the current context", func() {
				Expect(context).NotTo(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetCurrentContextNetwork", func() {
		var (
			network *environment.Network
			err     error
		)
		JustBeforeEach(func() {
			network, err = config.GetCurrentContextNetwork()
		})

		It("should fail without a current context", func() {
			Expect(network).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		Context("when current network is not set", func() {
			BeforeEach(func() {
				config = &environment.Config{
					CurrentContext: "foo",
					Contexts: map[string]*environment.Context{
						"foo": {},
					},
				}
			})

			It("should fail to return current network", func() {
				Expect(network).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when current network doesnt not exist", func() {
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

			It("should fail to return current network", func() {
				Expect(network).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when current network is set and exits", func() {
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

			It("should return the current network", func() {
				Expect(network).NotTo(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
})
