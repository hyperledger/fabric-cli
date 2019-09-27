/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
	"github.com/hyperledger/fabric-cli/pkg/plugin/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func TestFabricCommand(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Fabric Suite")
}

var _ = Describe("DefaultFabricCommand", func() {
	var (
		cmd      *cobra.Command
		settings *environment.Settings
		out      *bytes.Buffer
		err      *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)

		settings = environment.NewDefaultSettings()
		settings.Home = environment.Home(os.TempDir())
		settings.Streams = environment.Streams{
			Out: out,
			Err: err,
		}
	})

	JustBeforeEach(func() {
		cmd = NewDefaultFabricCommand(settings, []string{})
	})

	It("should create a fabric command", func() {
		Expect(cmd.Name()).To(Equal("fabric"))
		Expect(cmd.HasSubCommands()).To(BeTrue())
		Expect(cmd.Execute()).Should(Succeed())
		Expect(fmt.Sprint(out)).To(ContainSubstring("fabric [command]"))
		Expect(fmt.Sprint(out)).To(ContainSubstring("network"))
		Expect(fmt.Sprint(out)).To(ContainSubstring("context"))
		Expect(fmt.Sprint(out)).To(ContainSubstring("channel"))
		Expect(fmt.Sprint(out)).To(ContainSubstring("plugin"))
		Expect(fmt.Sprint(out)).To(ContainSubstring("chaincode"))
	})
})

var _ = Describe("LoadPlugins", func() {
	var (
		err      error
		cmd      *cobra.Command
		out      *bytes.Buffer
		settings *environment.Settings
		handler  *mocks.PluginHandler
	)

	BeforeEach(func() {
		cmd = &cobra.Command{}
		out = new(bytes.Buffer)

		cmd.SetOutput(out)

		settings = &environment.Settings{
			Home: environment.Home(os.TempDir()),
			Streams: environment.Streams{
				Out: out,
			},
		}

		handler = &mocks.PluginHandler{}
	})

	JustBeforeEach(func() {
		err = loadPlugins(cmd, settings, handler)
	})

	Context("when plugins are disabled", func() {
		BeforeEach(func() {
			settings.DisablePlugins = true

			err := handler.InstallPlugin("foo")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should not have foo command", func() {
			Expect(cmd.Execute()).Should(Succeed())
			Expect(cmd.HasSubCommands()).Should(BeFalse())
		})
	})

	Context("when plugin handler fails", func() {
		BeforeEach(func() {
			handler.GetPluginsReturns(nil, errors.New("handler error"))
			handler.LoadGoPluginReturns(nil, plugin.ErrNotAGoPlugin)
		})

		It("should fail to load plugins", func() {
			Expect(err.Error()).To(ContainSubstring("handler error"))
		})
	})

	Context("when plugins have been installed", func() {
		BeforeEach(func() {
			handler.GetPluginsReturns([]*plugin.Plugin{
				{
					Name: "foo",
					Command: &plugin.Command{
						Base: "./plugins",
					},
				},
			}, nil)
			handler.LoadGoPluginReturns(nil, plugin.ErrNotAGoPlugin)
		})

		It("should load plugins", func() {
			Expect(cmd.Execute()).Should(Succeed())
			Expect(cmd.HasSubCommands()).Should(BeTrue())
		})
	})

	Context("when Go plugins have been installed", func() {
		BeforeEach(func() {
			handler.GetPluginsReturns([]*plugin.Plugin{
				{
					Name: "foo",
					Command: &plugin.Command{
						Base: "./plugins",
					},
				},
			}, nil)
			handler.LoadGoPluginReturns(&cobra.Command{
				Run: func(cmd *cobra.Command, args []string) {},
			}, nil)
		})

		It("should load Go plugins", func() {
			Expect(cmd.Execute()).Should(Succeed())
			Expect(cmd.HasSubCommands()).Should(BeTrue())
		})
	})

	Context("when invalid Go plugins have been installed", func() {
		BeforeEach(func() {
			handler.GetPluginsReturns([]*plugin.Plugin{
				{
					Name: "foo",
					Command: &plugin.Command{
						Base: "./plugins",
					},
				},
			}, nil)
			handler.LoadGoPluginReturns(nil, fmt.Errorf("some Go plugin error"))
		})

		It("should fail to load Go plugins", func() {
			Expect(err.Error()).To(ContainSubstring("some Go plugin error"))
		})
	})
})
