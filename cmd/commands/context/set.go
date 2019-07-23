/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewContextSetCommand creates a new "fabric context set" command
func NewContextSetCommand(settings *environment.Settings) *cobra.Command {
	c := SetCommand{
		Context: new(environment.Context),
	}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "set [context-name]",
		Short: "set a context",
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
	flags.StringVar(&c.Context.Network, "network", "", "set the network context")
	flags.StringVar(&c.Context.Organization, "organization", "", "set the organization context")
	flags.StringVar(&c.Context.Channel, "channel", "", "set the channel context")
	flags.StringVar(&c.Context.User, "user", "", "set the users context")
	flags.StringArrayVar(&c.Context.Peers, "peers", []string{}, "set the peers context")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// SetCommand implements the set context command
type SetCommand struct {
	common.Command

	Name    string
	Context *environment.Context
}

// Validate checks the required parameters for run
func (c *SetCommand) Validate() error {
	if c.Context == nil ||
		(len(c.Context.Network) == 0 &&
			len(c.Context.Organization) == 0 &&
			len(c.Context.User) == 0) {
		return errors.New("context details not specified")
	}

	return nil
}

// Run executes the command
func (c *SetCommand) Run() error {
	if len(c.Name) == 0 {
		// ensure current context is set and exists
		if _, err := c.Settings.Config.GetCurrentContext(); err != nil {
			return err
		}

		c.Name = c.Settings.Config.CurrentContext
	}

	err := c.Settings.ModifyConfig(environment.SetContext(c.Name, c.Context))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully set context '%s'\n", c.Name)
	fmt.Fprintln(c.Settings.Streams.Out, "")
	fmt.Fprintln(c.Settings.Streams.Out, c.Settings.Config.Contexts[c.Name])

	return nil
}
