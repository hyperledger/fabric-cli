/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewPluginInstallCommand creates a new "fabric plugin install" command
func NewPluginInstallCommand(settings *environment.Settings) *cobra.Command {
	c := InstallCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "install <plugin-path>",
		Short: "Install a plugin from the local filesystem",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := c.Complete(); err != nil {
				return err
			}

			return c.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Path)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// InstallCommand implements the plugin install command
type InstallCommand struct {
	common.Command
	Handler plugin.Handler

	Path string
}

// Validate checks the required parameters for run
func (c *InstallCommand) Validate() error {
	if len(c.Path) == 0 {
		return errors.New("plugin path not specified")
	}

	return nil
}

// Run executes the command
func (c *InstallCommand) Run() error {
	err := c.Handler.InstallPlugin(c.Path)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, "successfully installed the plugin")

	return nil
}

// Complete initializes the plugin handler
func (c *InstallCommand) Complete() error {
	c.Handler = &plugin.DefaultHandler{
		Dir:      c.Settings.Home.Plugins(),
		Filename: plugin.DefaultFilename,
	}
	return nil
}
