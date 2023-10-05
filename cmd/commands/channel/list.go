/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-admin-sdk/pkg/channel"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelListCommand creates a new "fabric channel list" command
func NewChannelListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all joined channels",
		Long:  "List all joined channels, peer is the current context's peer",
		/*PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Complete()
		},*/
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
	context := context.Background()
	peerChannelInfo, err := channel.ListChannelOnPeer(context, c.Connection, c.OrgMSP)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, "Channels Joined:")
	for _, channel := range peerChannelInfo {
		fmt.Fprintf(c.Settings.Streams.Out, " - %s\n", channel.ChannelId)
	}

	return nil
}
