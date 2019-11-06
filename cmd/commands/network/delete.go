/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package network

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewNetworkDeleteCommand creates a new "fabric network delete" command
func NewNetworkDeleteCommand(settings *environment.Settings) *cobra.Command {
	c := DeleteCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "delete <network-name>",
		Short: "delete a network",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Name)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// DeleteCommand implements the network delete command
type DeleteCommand struct {
	common.Command

	Name string
}

// Validate checks the required parameters for run
func (c *DeleteCommand) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("network name not specified")
	}

	return nil
}

// Run executes the command
func (c *DeleteCommand) Run() error {

	if _, ok := c.Settings.Config.Networks[c.Name]; !ok {
		err := fmt.Sprintf("network %s doesn't exist", c.Name)
		return errors.New(err)
	}

	err := c.Settings.ModifyConfig(environment.DeleteNetwork(c.Name))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully deleted network '%s'\n", c.Name)

	return nil
}
