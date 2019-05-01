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

// NewNetworkSetCommand creates a new "fabric network set" command
func NewNetworkSetCommand(settings *environment.Settings) *cobra.Command {
	c := SetCommand{
		Network: new(environment.Network),
	}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "set <network-name>",
		Short: "set a network",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Name)

	flags := cmd.Flags()
	flags.StringVar(&c.Network.ConfigPath, "path", "", "set the path to the network configurations")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// SetCommand implements the network set command
type SetCommand struct {
	common.Command

	Name    string
	Network *environment.Network
}

// Validate checks the required parameters for run
func (c *SetCommand) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("network name not specified")
	}

	if c.Network == nil || len(c.Network.ConfigPath) == 0 {
		return errors.New("network configuration path not specified")
	}

	return nil
}

// Run executes the command
func (c *SetCommand) Run() error {
	err := c.Settings.ModifyConfig(environment.SetNetwork(c.Name, c.Network))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully set network '%s'\n", c.Name)
	fmt.Fprintln(c.Settings.Streams.Out, "")
	fmt.Fprintln(c.Settings.Streams.Out, c.Settings.Config.Networks[c.Name])

	return nil
}
