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

// NewPluginInstallCommand creates a new "fabric plugin install" command
func NewPluginInstallCommand(settings *environment.Settings) *cobra.Command {
	c := InstallCommand{
		Out: settings.Streams.Out,
		Handler: &plugin.DefaultHandler{
			Dir:      settings.Home.Plugins(),
			Filename: plugin.DefaultFilename,
		},
	}

	cmd := &cobra.Command{
		Use:   "install <path>",
		Short: "Install a plugin from the local filesystem",
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

// InstallCommand implements the plugin install command
type InstallCommand struct {
	Out     io.Writer
	Handler plugin.Handler

	path string
}

// Complete populates required fields for Run
func (cmd *InstallCommand) Complete(args []string) error {
	if len(args) == 0 {
		return errors.New("plugin path not specified")
	}

	cmd.path = strings.TrimSpace(args[0])

	if len(cmd.path) == 0 {
		return errors.New("plugin path not specified")
	}

	return nil
}

// Run executes the command
func (cmd *InstallCommand) Run() error {
	err := cmd.Handler.InstallPlugin(cmd.path)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.Out, "successfully installed the plugin")

	return nil
}
