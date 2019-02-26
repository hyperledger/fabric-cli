/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewPluginCommand creates a new "fabric plugin" command
func NewPluginCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
	}

	cmd.AddCommand(
		NewPluginListCommand(settings),
		NewPluginInstallCommand(settings),
		NewPluginUninstallCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}
