/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelListCommand creates a new "fabric channel list" command
func NewChannelListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all joined channels",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Complete()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ListCommand implements the channel list command
type ListCommand struct {
	BaseCommand
}

// Run executes the command
func (c *ListCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
	}

	resp, err := c.ResourceManagement.QueryChannels(options...)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, "Channels Joined:")
	for _, channel := range resp.Channels {
		fmt.Fprintf(c.Settings.Streams.Out, " - %s\n", channel.ChannelId)
	}

	return nil
}
