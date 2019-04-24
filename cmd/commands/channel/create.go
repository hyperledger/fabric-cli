/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

// NewChannelCreateCommand creates a new "fabric channel create" command
func NewChannelCreateCommand(settings *environment.Settings) *cobra.Command {
	c := CreateCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "create <channel-id> <tx-path>",
		Short: "create a new channel",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := c.Complete(cmd); err != nil {
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

	cmd.SetOutput(c.Out)

	return cmd
}

// CreateCommand implements the channel create command
type CreateCommand struct {
	Out      io.Writer
	Settings *environment.Settings
	Profile  *environment.Profile

	ResourceManangement fabric.ResourceManagement

	ChannelID string
	ChannelTX string
}

// Complete populates required fields for Run
func (c *CreateCommand) Complete(cmd *cobra.Command) error {
	var err error

	c.Profile, err = c.Settings.GetActiveProfile()
	if err != nil {
		return err
	}

	if c.ResourceManangement == nil {
		c.ResourceManangement, err = fabric.NewResourceManagementClient(c.Profile)
		if err != nil {
			return err
		}
	}

	if c.Profile.Context != nil {
		c.ChannelID = c.Profile.Context.Channel
	}

	args := cmd.Flags().Args()
	if len(args) != 2 {
		return fmt.Errorf("unexpected args: %v", args)
	}

	c.ChannelID = strings.TrimSpace(args[0])
	c.ChannelTX = strings.TrimSpace(args[1])

	return nil
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

	if _, err := c.ResourceManangement.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     c.ChannelID,
		ChannelConfig: r,
	}); err != nil {
		return err
	}

	fmt.Fprintf(c.Out, "successfully created channel '%s'\n", c.ChannelID)

	return nil
}
