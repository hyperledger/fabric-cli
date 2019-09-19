/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewPluginListCommand creates a new "fabric plugin list" command
func NewPluginListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all installed plugins",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.Complete()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ListCommand implements the plugin list command
type ListCommand struct {
	common.Command
	Handler plugin.Handler
}

// Run executes the command
func (c *ListCommand) Run() error {
	plugins, err := c.Handler.GetPlugins()
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		fmt.Fprintln(c.Settings.Streams.Out, "no plugins currently exist")
		return nil
	}

	for _, plugin := range plugins {
		fmt.Fprint(c.Settings.Streams.Out, plugin.Name, "\n")
	}

	return nil
}

// Complete initializes the plugin handler
func (c *ListCommand) Complete() error {
	c.Handler = &plugin.DefaultHandler{
		Dir:      c.Settings.Home.Plugins(),
		Filename: plugin.DefaultFilename,
	}
	return nil
}
