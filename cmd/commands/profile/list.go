/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile

import (
	"errors"
	"fmt"
	"io"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewProfileListCommand creates a new "fabric profile list" command
func NewProfileListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{
		Out:      settings.Streams.Out,
		Profiles: settings.Profiles,
		Active:   settings.ActiveProfile,
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
	Profiles []*environment.Profile
	Active   string
}

// Run executes the command
func (cmd *ListCommand) Run() error {
	if len(cmd.Profiles) == 0 {
		return errors.New("no profiles currently exist")
	}

	for _, p := range cmd.Profiles {
		fmt.Fprint(cmd.Out, p.Name)
		if p.Name == cmd.Active {
			fmt.Fprint(cmd.Out, " (active)")
		}
		fmt.Fprintln(cmd.Out)

	}

	return nil
}
