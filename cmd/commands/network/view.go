/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network

import (
	"fmt"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewNetworkViewCommand creates a new "fabric network view" command
func NewNetworkViewCommand(settings *environment.Settings) *cobra.Command {
	c := ViewCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "view [network-name]",
		Short: "view a network",
		Args:  c.ParseArgs(),
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Name)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ViewCommand implements the current command
type ViewCommand struct {
	common.Command

	Name string
}

// Run executes the command
func (c *ViewCommand) Run() error {
	if len(c.Name) == 0 {
		context, err := c.Settings.Config.GetCurrentContext()
		if err != nil {
			return err
		}

		c.Name = context.Network
	}

	network, ok := c.Settings.Config.Networks[c.Name]
	if !ok {
		return fmt.Errorf("network '%s' does not exist", c.Name)
	}

	fmt.Fprintln(c.Settings.Streams.Out, "Name:	", c.Name)
	fmt.Fprintln(c.Settings.Streams.Out, network)

	return nil
}
