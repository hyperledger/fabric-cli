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

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

// NewChannelConfigCommand creates a new "fabric channel config" command
func NewChannelConfigCommand(settings *environment.Settings) *cobra.Command {
	c := ConfigCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "config <channel-id>",
		Short: "get the channel configuration",
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

// ConfigCommand implements the channel config command
type ConfigCommand struct {
	Out      io.Writer
	Settings *environment.Settings
	Profile  *environment.Profile

	ResourceManangement fabric.ResourceManagement

	ChannelID string
}

// Complete populates required fields for Run
func (c *ConfigCommand) Complete(cmd *cobra.Command) error {
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

	args := cmd.Flags().Args()
	if len(args) != 1 {
		return fmt.Errorf("unexpected args: %v", args)
	}

	c.ChannelID = strings.TrimSpace(args[0])

	return nil
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
	resp, err := c.ResourceManangement.QueryConfigFromOrderer(c.ChannelID)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Out, "ID: %s\n", resp.ID())
	fmt.Fprintf(c.Out, "Latest Block Number: %d\n", resp.BlockNumber())

	if len(resp.Orderers()) > 0 {
		fmt.Fprintln(c.Out, "Orderers:")
		for _, orderer := range resp.Orderers() {
			fmt.Fprintf(c.Out, " - %s\n", orderer)
		}
	}

	if len(resp.AnchorPeers()) > 0 {
		fmt.Fprintln(c.Out, "Anchor Peers:")
		for _, anchor := range resp.AnchorPeers() {
			fmt.Fprintf(c.Out, " - %s:%d (%s)\n", anchor.Host, anchor.Port, anchor.Org)
		}
	}

	return nil
}
