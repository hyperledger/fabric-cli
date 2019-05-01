/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	"fmt"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewContextViewCommand creates a new "fabric context view" command
func NewContextViewCommand(settings *environment.Settings) *cobra.Command {
	c := ViewCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "view [context-name]",
		Short: "view a context",
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
		c.Name = c.Settings.Config.CurrentContext
	}

	context, ok := c.Settings.Config.Contexts[c.Name]
	if !ok {
		return fmt.Errorf("context '%s' does not exist", c.Name)
	}

	fmt.Fprintln(c.Settings.Streams.Out, "Name:		", c.Name)
	fmt.Fprintln(c.Settings.Streams.Out, context)

	return nil
}
