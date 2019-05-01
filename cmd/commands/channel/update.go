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

// NewChannelUpdateCommand creates a new "fabric channel update" command
func NewChannelUpdateCommand(settings *environment.Settings) *cobra.Command {
	c := UpdateCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "update <channel-id> <tx-path>",
		Short: "update a channel",
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

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// UpdateCommand implements the channel update command
type UpdateCommand struct {
	BaseCommand

	ChannelID string
	ChannelTX string
}

// Validate checks the required parameters for run
func (c *UpdateCommand) Validate() error {
	if len(c.ChannelID) == 0 {
		return errors.New("channel id not specified")
	}

	if len(c.ChannelTX) == 0 {
		return errors.New("channel tx path not specified")
	}

	return nil
}

// Run executes the command
func (c *UpdateCommand) Run() error {
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

	fmt.Fprintf(c.Settings.Streams.Out, "successfully updated channel '%s'\n", c.ChannelID)

	return nil
}
