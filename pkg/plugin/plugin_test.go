/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlugin(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Plugin Suite")
}

//go:generate gobin -m -run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/handler.go --fake-name PluginHandler . Handler

var _ = Describe("Plugin", func() {
	var (
		dir     string
		handler *plugin.DefaultHandler
	)

	BeforeEach(func() {
		dir = filepath.Join(os.TempDir(), "plugins")

		handler = &plugin.DefaultHandler{
			Dir:      dir,
			Filename: plugin.DefaultFilename,
		}
	})

	JustBeforeEach(func() {
		err := os.RemoveAll(handler.Dir)

		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(handler.Dir)

		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Get", func() {
		It("should get plugins", func() {
			plugins, err := handler.GetPlugins()

			Expect(err).To(BeNil())
			Expect(len(plugins)).To(Equal(0))
		})

		Context("when the plugin directory is malformed", func() {
			BeforeEach(func() {
				handler.Dir = "[]"
			})

			It("should fail to get plugins", func() {
				_, err := handler.GetPlugins()

				Expect(err).NotTo(BeNil())
			})
		})

		Context("when the plugin directory is empty", func() {
			BeforeEach(func() {
				handler.Dir = ""
			})

			It("should fail to get plugins", func() {
				_, err := handler.GetPlugins()

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Install", func() {
		It("should fail to uninstall non-existent plugin", func() {
			err := handler.UninstallPlugin("foo")

			Expect(err).NotTo(BeNil())
		})

		Context("when plugin yaml does not exist", func() {
			BeforeEach(func() {
				err := os.MkdirAll(filepath.Join(os.TempDir(), "foo", "bar"), 0777)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail to install plugin", func() {
				err := handler.InstallPlugin(filepath.Join(os.TempDir(), "foo", "bar"))

				Expect(err).NotTo(BeNil())
			})
		})

		Context("when plugin yaml is malformed", func() {
			BeforeEach(func() {
				err := os.MkdirAll(filepath.Join(os.TempDir(), "foo", "bar"), 0777)

				Expect(err).NotTo(HaveOccurred())

				err = ioutil.WriteFile(filepath.Join(os.TempDir(), "foo", "bar", handler.Filename),
					[]byte("command: !!float 'error'"), 0777)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail to install plugin", func() {
				err := handler.InstallPlugin(filepath.Join(os.TempDir(), "foo", "bar"))

				Expect(err).NotTo(BeNil())
			})
		})

		Context("when a plugin is installed", func() {
			JustBeforeEach(func() {
				err := handler.InstallPlugin("./testdata/plugins/home")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail to install existing plugin", func() {
				err := handler.InstallPlugin("./testdata/plugins/home")

				Expect(err).NotTo(BeNil())

				plugins, err := handler.GetPlugins()

				Expect(err).To(BeNil())
				Expect(len(plugins)).To(Equal(1))
			})
		})
	})

	Describe("Uninstall", func() {
		It("should install plugin", func() {
			Expect(handler.InstallPlugin("./testdata/plugins/home")).Should(Succeed())

			plugins, err := handler.GetPlugins()

			Expect(err).To(BeNil())
			Expect(len(plugins)).To(Equal(1))
		})

		It("should fail to install non-existent plugin", func() {
			err := handler.InstallPlugin(".foo/bar")

			Expect(err).NotTo(BeNil())
		})

		Context("when a plugin is installed", func() {
			JustBeforeEach(func() {
				err := handler.InstallPlugin("./testdata/plugins/home")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should uninstall plugin", func() {
				Expect(handler.UninstallPlugin("home")).Should(Succeed())

				plugins, err := handler.GetPlugins()

				Expect(err).To(BeNil())
				Expect(len(plugins)).To(Equal(0))
			})
		})
	})

	Describe("LoadGoPlugin", func() {
		It("should load Go plugin", func() {
			tmpdir, err := ioutil.TempDir("", "echogoplugin")
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				err := os.RemoveAll(tmpdir)
				Expect(err).NotTo(HaveOccurred())
			}()

			pluginPath := filepath.Join(tmpdir, "echogoplugin")

			err = buildGoPlugin("./testdata/plugins/echogoplugin/cmd/echogoplugin.go", pluginPath)
			Expect(err).NotTo(HaveOccurred())

			cmd, err := handler.LoadGoPlugin(pluginPath, &environment.Settings{Streams: environment.Streams{Out: os.Stdout}})
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).NotTo(BeNil())
			cmd.SetArgs([]string{"--message", "Hello world!"})
			Expect(cmd.Execute()).NotTo(HaveOccurred())
		})
	})

})

func buildGoPlugin(path, outputPath string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, path)
	_, err := cmd.CombinedOutput()
	return err
}
