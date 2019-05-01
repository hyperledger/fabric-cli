/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChannelConfigCommand creates a new "fabric channel config" command
func NewChannelConfigCommand(settings *environment.Settings) *cobra.Command {
	c := ConfigCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "config <channel-id>",
		Short: "get the channel configuration",
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

// ConfigCommand implements the channel config command
type ConfigCommand struct {
	BaseCommand

	ChannelID string
}

// Validate checks the required parameters for run
func (c *ConfigCommand) Validate() error {
	if len(c.ChannelID) == 0 {
		return errors.New("channel id not specified")
	}

	return nil
}

// Run executes the command
func (c *ConfigCommand) Run() error {
	resp, err := c.ResourceManagement.QueryConfigFromOrderer(c.ChannelID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "ID: %s\n", resp.ID())
	fmt.Fprintf(c.Settings.Streams.Out, "Latest Block Number: %d\n", resp.BlockNumber())

	if len(resp.Orderers()) > 0 {
		fmt.Fprintln(c.Settings.Streams.Out, "Orderers:")
		for _, orderer := range resp.Orderers() {
			fmt.Fprintf(c.Settings.Streams.Out, " - %s\n", orderer)
		}
	}

	if len(resp.AnchorPeers()) > 0 {
		fmt.Fprintln(c.Settings.Streams.Out, "Anchor Peers:")
		for _, anchor := range resp.AnchorPeers() {
			fmt.Fprintf(c.Settings.Streams.Out, " - %s:%d (%s)\n", anchor.Host, anchor.Port, anchor.Org)
		}
	}

	return nil
}
