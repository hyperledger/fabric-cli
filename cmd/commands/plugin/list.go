/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

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
	pcmd := pluginListCommand{
		out: settings.Streams.Out,
		handler: &plugin.DefaultHandler{
			Dir:              settings.Home.Plugins(),
			MetadataFileName: plugin.DefaultMetadataFileName,
		},
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type pluginListCommand struct {
	out     io.Writer
	handler plugin.Handler
}

func (cmd *pluginListCommand) run() error {
	plugins, err := cmd.handler.GetPlugins()
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		fmt.Fprint(cmd.out, plugin.Name, "\n")
	}

	return nil
}
