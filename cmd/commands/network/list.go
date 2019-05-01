/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewNetworkListCommand creates a new "fabric network list" command
func NewNetworkListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all networks",
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ListCommand implements the network list command
type ListCommand struct {
	common.Command
}

// Run executes the command
func (c *ListCommand) Run() error {
	if len(c.Settings.Config.Networks) == 0 {
		return errors.New("no networks currently exist")
	}

	context, _ := c.Settings.Config.GetCurrentContext()

	var names []string
	for name := range c.Settings.Config.Networks {
		if context != nil && name == context.Network {
			name += " (current)"
		}
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintln(c.Settings.Streams.Out, name)
	}

	return nil
}
