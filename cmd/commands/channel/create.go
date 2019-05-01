/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelCreateCommand creates a new "fabric channel create" command
func NewChannelCreateCommand(settings *environment.Settings) *cobra.Command {
	c := &CreateCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "create <channel-id> <tx-path>",
		Short: "create a new channel",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, args []string) error {
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
	c.AddArg(&c.ChannelTX)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// CreateCommand implements the channel create command
type CreateCommand struct {
	BaseCommand

	ChannelID string
	ChannelTX string
}

// Validate checks the required parameters for run
func (c *CreateCommand) Validate() error {
	if len(c.ChannelID) == 0 {
		return errors.New("channel id not specified")
	}

	if len(c.ChannelTX) == 0 {
		return errors.New("channel tx path not specified")
	}

	return nil
}

// Run executes the command
func (c *CreateCommand) Run() error {
	r, err := os.Open(c.ChannelTX)
	if err != nil {
		return err
	}

	defer r.Close()

	if _, err := c.ResourceManagement.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     c.ChannelID,
		ChannelConfig: r,
	}); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully created channel '%s'\n", c.ChannelID)

	return nil
}
