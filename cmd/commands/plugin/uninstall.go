/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewPluginUninstallCommand creates a new "fabric plugin uninstall" command
func NewPluginUninstallCommand(settings *environment.Settings) *cobra.Command {
	c := UninstallCommand{
		Out: settings.Streams.Out,
		Handler: &plugin.DefaultHandler{
			Dir:      settings.Home.Plugins(),
			Filename: plugin.DefaultFilename,
		},
	}

	cmd := &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Uninstall a plugin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.Complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Out)

	return cmd
}

// UninstallCommand implements the plugin uninstall command
type UninstallCommand struct {
	Out     io.Writer
	Handler plugin.Handler

	name string
}

// Complete populates required fields for Run
func (cmd *UninstallCommand) Complete(args []string) error {
	if len(args) == 0 {
		return errors.New("plugin name not specified")
	}

	cmd.name = strings.TrimSpace(args[0])

	if len(cmd.name) == 0 {
		return errors.New("plugin name not specified")
	}

	return nil
}

// Run executes the command
func (cmd *UninstallCommand) Run() error {
	err := cmd.Handler.UninstallPlugin(cmd.name)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.Out, "successfully uninstalled the plugin")

	return nil
}
