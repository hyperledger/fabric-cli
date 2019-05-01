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

// NewContextDeleteCommand creates a new "fabric context delete" command
func NewContextDeleteCommand(settings *environment.Settings) *cobra.Command {
	c := DeleteCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "delete <context-name>",
		Short: "delete a context",
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

// DeleteCommand implements the delete context command
type DeleteCommand struct {
	common.Command

	Name string
}

// Validate checks the required parameters for run
func (c *DeleteCommand) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("context name not specified")
	}

	return nil
}

// Run executes the command
func (c *DeleteCommand) Run() error {
	err := c.Settings.ModifyConfig(environment.DeleteContext(c.Name))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully deleted context '%s'\n", c.Name)

	return nil
}
