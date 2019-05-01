/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelJoinCommand creates a new "fabric channel join" command
func NewChannelJoinCommand(settings *environment.Settings) *cobra.Command {
	c := JoinCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "join <channel-id>",
		Short: "join a channel",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := c.Complete(); err != nil {
				return err
			}

			if err := c.Validate(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.ChannelID)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// JoinCommand implements the channel join command
type JoinCommand struct {
	BaseCommand

	ChannelID string
}

// Validate checks the required parameters for run
func (c *JoinCommand) Validate() error {
	if len(c.ChannelID) == 0 {
		return errors.New("channel id not specified")
	}

	return nil
}

// Run executes the command
func (c *JoinCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
	}

	if err := c.ResourceManagement.JoinChannel(c.ChannelID, options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully joined channel '%s'\n", c.ChannelID)

	return nil
}
