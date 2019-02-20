/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewPluginUninstallCommand creates a new "fabric plugin uninstall" command
func NewPluginUninstallCommand(settings *environment.Settings) *cobra.Command {
	pcmd := pluginUninstallCommand{
		out: settings.Streams.Out,
		handler: &plugin.DefaultHandler{
			Dir:      settings.Home.Plugins(),
			Filename: plugin.DefaultFilename,
		},
	}

	cmd := &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Uninstall a plugin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type pluginUninstallCommand struct {
	out     io.Writer
	handler plugin.Handler

	name string
}

func (cmd *pluginUninstallCommand) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("plugin name not specified")
	}

	cmd.name = args[0]

	return nil
}

func (cmd *pluginUninstallCommand) run() error {
	err := cmd.handler.UninstallPlugin(cmd.name)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.out, "successfully uninstalled the plugin")

	return nil
}
