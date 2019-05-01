/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network

import (
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewNetworkCommand creates a new "fabric network" command
func NewNetworkCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage networks",
	}

	cmd.AddCommand(
		NewNetworkViewCommand(settings),
		NewNetworkListCommand(settings),
		NewNetworkSetCommand(settings),
		NewNetworkDeleteCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}
