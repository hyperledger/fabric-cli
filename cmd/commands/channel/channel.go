/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelCommand creates a new "fabric channel" command
func NewChannelCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Manage channels",
	}

	cmd.AddCommand(
		NewChannelCreateCommand(settings),
		NewChannelJoinCommand(settings),
		NewChannelUpdateCommand(settings),
		NewChannelListCommand(settings),
		NewChannelConfigCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}
