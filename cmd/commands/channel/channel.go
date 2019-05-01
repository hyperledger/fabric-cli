/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
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

// BaseCommand implements common channel command functions
type BaseCommand struct {
	common.Command

	Factory            fabric.Factory
	ResourceManagement fabric.ResourceManagement
}

// Complete initializes all clients needed for Run
func (c *BaseCommand) Complete() error {
	var err error

	if c.Factory == nil {
		c.Factory, err = fabric.NewFactory(c.Settings.Config)
		if err != nil {
			return err
		}
	}

	c.ResourceManagement, err = c.Factory.ResourceManagement()
	if err != nil {
		return err
	}

	return nil
}
