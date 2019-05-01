/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewContextListCommand creates a new "fabric context list" command
func NewContextListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all contexts",
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ListCommand implements the list context command
type ListCommand struct {
	common.Command
}

// Run executes the command
func (c *ListCommand) Run() error {
	if len(c.Settings.Config.Contexts) == 0 {
		return errors.New("no contexts currently exist")
	}

	var names []string
	for name := range c.Settings.Config.Contexts {
		if name == c.Settings.Config.CurrentContext {
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
