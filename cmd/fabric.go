/*
Copyright State Street Corp. All Rights Reserved.

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

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}

// NewDefaultFabricCommand returns a new default root commad for fabric
func NewDefaultFabricCommand(settings *environment.Settings, args []string) *cobra.Command {
	cmd := NewFabricCommand(settings)
	flags := cmd.PersistentFlags()

	settings.AddFlags(flags)
	flags.Parse(args)

	if err := settings.Init(flags); err != nil {
		fmt.Fprintf(settings.Streams.Err, "An error occurred while loading configurations: %v\n", err)
		os.Exit(1)
	}

	// must resolve home and config file before loading config flags
	settings.Config.AddFlags(flags)

	if err := loadPlugins(cmd, settings, &plugin.DefaultHandler{
		Dir:      settings.Home.Plugins(),
		Filename: plugin.DefaultFilename,
	}); err != nil {
		fmt.Fprintf(settings.Streams.Err, "An error occurred while loading plugins: %v\n", err)
		os.Exit(1)
	}

	return cmd
}

func main() {
	settings := environment.NewDefaultSettings()
	cmd := NewDefaultFabricCommand(settings, os.Args[1:])

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

// loadPlugins processes all of the installed plugins, wraps them with cobra,
// and adds them to the root command
func loadPlugins(cmd *cobra.Command, settings *environment.Settings, handler plugin.Handler) error {
	if settings.DisablePlugins {
		return nil
	}

	plugins, err := handler.GetPlugins()
	if err != nil {
		return err
	}

	settings.SetupPluginEnvironment()

	for _, p := range plugins {
		c, err := loadPlugin(p, settings, handler)
		if err != nil {
			return err
		}
		cmd.AddCommand(c)
	}

	return nil
}

// loadPlugin loads the given plugin as either a Go plugin or a wrapped executable
func loadPlugin(p *plugin.Plugin, settings *environment.Settings, handler plugin.Handler) (*cobra.Command, error) {
	path := os.ExpandEnv(p.Command.Base)
	c, err := handler.LoadGoPlugin(path, settings)
	if err == nil {
		return c, nil
	}
	if err != plugin.ErrNotAGoPlugin {
		return nil, err
	}

	return &cobra.Command{
		Use:   p.Name,
		Short: p.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			e := exec.Command(path, append(p.Command.Args, args...)...)
			e.Env = os.Environ()
			e.Stdin = settings.Streams.In
			e.Stdout = settings.Streams.Out
			e.Stderr = settings.Streams.Err
			return e.Run()
		},
	}, nil
}
