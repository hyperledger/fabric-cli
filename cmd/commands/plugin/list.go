/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewPluginListCommand creates a new "fabric plugin list" command
func NewPluginListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{
		Out: settings.Streams.Out,
		Handler: &plugin.DefaultHandler{
			Dir:      settings.Home.Plugins(),
			Filename: plugin.DefaultFilename,
		},
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Out)

	return cmd
}

// ListCommand implements the plugin list command
type ListCommand struct {
	Out     io.Writer
	Handler plugin.Handler
}

// Run executes the command
func (cmd *ListCommand) Run() error {
	plugins, err := cmd.Handler.GetPlugins()
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		fmt.Fprintln(cmd.Out, "no plugins currently exist")
		return nil
	}

	for _, plugin := range plugins {
		fmt.Fprint(cmd.Out, plugin.Name, "\n")
	}

	return nil
}
