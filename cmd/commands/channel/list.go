/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"
	"io"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

// NewChannelListCommand creates a new "fabric channel list" command
func NewChannelListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all joined channels",
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

// ListCommand implements the channel list command
type ListCommand struct {
	Out      io.Writer
	Settings *environment.Settings
	Profile  *environment.Profile

	ResourceManangement fabric.ResourceManagement
	Options             []resmgmt.RequestOption
}

// Complete populates required fields for Run
func (c *ListCommand) Complete(cmd *cobra.Command) error {
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

	// TODO: support multiple peers with queries in parallel
	if c.Options == nil && c.Profile.Context != nil && len(c.Profile.Context.Peers) == 1 {
		c.Options = []resmgmt.RequestOption{
			resmgmt.WithTargetEndpoints(c.Profile.Peers[c.Profile.Context.Peers[0]].URL),
		}
	}

	return nil
}

// Validate checks the required parameters for run
func (c *ListCommand) Validate() error {
	if len(c.Options) == 0 {
		return errors.New("peer not specified")
	}

	return nil
}

// Run executes the command
func (c *ListCommand) Run() error {
	resp, err := c.ResourceManangement.QueryChannels(c.Options...)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Out, "Channels Joined:")
	for _, channel := range resp.Channels {
		fmt.Fprintf(c.Out, " - %s\n", channel.ChannelId)
	}

	return nil
}
