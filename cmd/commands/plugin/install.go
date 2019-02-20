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

// NewPluginInstallCommand creates a new "fabric plugin install" command
func NewPluginInstallCommand(settings *environment.Settings) *cobra.Command {
	pcmd := pluginInstallCommand{
		out: settings.Streams.Out,
		handler: &plugin.DefaultHandler{
			Dir:      settings.Home.Plugins(),
			Filename: plugin.DefaultFilename,
		},
	}

	cmd := &cobra.Command{
		Use:   "install <path>",
		Short: "Install a plugin from the local filesystem",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type pluginInstallCommand struct {
	out     io.Writer
	handler plugin.Handler

	path string
}

func (cmd *pluginInstallCommand) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("plugin path not specified")
	}

	cmd.path = args[0]

	return nil
}

func (cmd *pluginInstallCommand) run() error {
	err := cmd.handler.InstallPlugin(cmd.path)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.out, "successfully installed the plugin")

	return nil
}
