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

// NewPluginUninstallCommand creates a new "fabric plugin uninstall" command
func NewPluginUninstallCommand(settings *environment.Settings) *cobra.Command {
	c := UninstallCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "uninstall <plugin-name>",
		Short: "Uninstall a plugin",
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

	c.AddArg(&c.Name)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// UninstallCommand implements the plugin uninstall command
type UninstallCommand struct {
	common.Command
	Handler plugin.Handler

	Name string
}

// Validate checks the required parameters for run
func (c *UninstallCommand) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("plugin name not specified")
	}

	return nil
}

// Run executes the command
func (c *UninstallCommand) Run() error {
	err := c.Handler.UninstallPlugin(c.Name)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, "successfully uninstalled the plugin")

	return nil
}

// Complete initializes the plugin handler
func (c *UninstallCommand) Complete() error {
	c.Handler = &plugin.DefaultHandler{
		Dir:      c.Settings.Home.Plugins(),
		Filename: plugin.DefaultFilename,
	}
	return nil
}
