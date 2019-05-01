/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewContextCommand creates a new "fabric context" command
func NewContextCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage contexts",
	}

	cmd.AddCommand(
		NewContextViewCommand(settings),
		NewContextUseCommand(settings),
		NewContextListCommand(settings),
		NewContextSetCommand(settings),
		NewContextDeleteCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}
