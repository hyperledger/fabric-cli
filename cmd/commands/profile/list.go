/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewProfileListCommand creates a new "fabric profile list" command
func NewProfileListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all configuration profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Out)

	return cmd
}

// ListCommand implements the profile list command
type ListCommand struct {
	Out      io.Writer
	Settings *environment.Settings
}

// Run executes the command
func (cmd *ListCommand) Run() error {
	if len(cmd.Settings.Profiles) == 0 {
		return errors.New("no profiles currently exist")
	}

	var keys []string
	for k := range cmd.Settings.Profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p := cmd.Settings.Profiles[k]
		fmt.Fprint(cmd.Out, p.Name)
		if p.Name == cmd.Settings.ActiveProfile {
			fmt.Fprint(cmd.Out, " (active)")
		}

		fmt.Fprintln(cmd.Out)
	}

	return nil
}
