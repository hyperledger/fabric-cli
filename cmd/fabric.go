/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewFabricCommand returns a new root command for fabric
func NewFabricCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fabric",
		Short: "The command line interface for Hyperledger Fabric",
	}

	// load all built in commands into the root command
	cmd.AddCommand(commands.All(settings)...)

	// load all plugins into the root command
	loadPlugins(cmd, settings, &plugin.DefaultHandler{
		Dir:              settings.Home.Plugins(),
		MetadataFileName: plugin.DefaultMetadataFileName,
	})

	return cmd
}

func main() {
	settings, err := environment.GetSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	cmd := NewFabricCommand(settings)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// loadPlugins processes all of the installed plugins, wraps them with cobra,
// and adds them to the root command
func loadPlugins(cmd *cobra.Command, settings *environment.Settings, handler plugin.Handler) {
	if settings.DisablePlugins {
		return
	}

	plugins, err := handler.GetPlugins()
	if err != nil {
		fmt.Fprintf(settings.Streams.Err, "An error occurred while loading plugins: %s", err)
		return
	}

	for _, plugin := range plugins {
		p := plugin
		c := &cobra.Command{
			Use:   p.Name,
			Short: p.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				e := exec.Command(os.ExpandEnv(p.Command.Base),
					append(p.Command.Args, args...)...)
				e.Stdin = settings.Streams.In
				e.Stdout = settings.Streams.Out
				e.Stderr = settings.Streams.Err
				return e.Run()
			},
		}

		cmd.AddCommand(c)
	}
}
