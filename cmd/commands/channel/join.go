/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

// NewChannelJoinCommand creates a new "fabric channel join" command
func NewChannelJoinCommand(settings *environment.Settings) *cobra.Command {
	c := JoinCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "join <channel-id>",
		Short: "join a channel",
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

// JoinCommand implements the channel join command
type JoinCommand struct {
	Out      io.Writer
	Settings *environment.Settings
	Profile  *environment.Profile

	ResourceManangement fabric.ResourceManagement
	Options             []resmgmt.RequestOption

	ChannelID string
}

// Complete populates required fields for Run
func (c *JoinCommand) Complete(cmd *cobra.Command) error {
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

	if c.Options == nil && c.Profile.Context != nil {
		peers := []string{}
		for _, p := range c.Profile.Context.Peers {
			peers = append(peers, c.Profile.Peers[p].URL)
		}

		c.Options = []resmgmt.RequestOption{
			resmgmt.WithTargetEndpoints(peers...),
		}
	}

	args := cmd.Flags().Args()
	if len(args) != 1 {
		return fmt.Errorf("unexpected args: %v", args)
	}

	c.ChannelID = strings.TrimSpace(args[0])

	return nil
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
	if err := c.ResourceManangement.JoinChannel(c.ChannelID, c.Options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Out, "successfully joined channel '%s'\n", c.ChannelID)

	return nil
}
