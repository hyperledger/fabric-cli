/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/hyperledger/fabric-admin-sdk/pkg/identity"
	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelCommand creates a new "fabric channel" command
func NewChannelCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Manage channels",
		Long:  "Manage channels with config|create|join|list|update",
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

// BaseCommand implements common channel command functions
type BaseCommand struct {
	common.Command
	OrgMSP     identity.SigningIdentity
	Connection *grpc.ClientConn
}

// Complete initializes all clients needed for Run
func (c *BaseCommand) Complete() error {
	return nil
}
