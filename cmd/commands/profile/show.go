/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewProfileShowCommand creates a new "fabric profile show" command
func NewProfileShowCommand(settings *environment.Settings) *cobra.Command {
	c := ShowCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "show [profilename]",
		Short: "show the metadata of the active profile",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.Complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run()
		},
	}

	cmd.SetOutput(c.Out)

	return cmd
}

// ShowCommand implements the profile show command
type ShowCommand struct {
	Out      io.Writer
	Settings *environment.Settings

	name string
}

// Complete populates required fields for Run
func (cmd *ShowCommand) Complete(args []string) error {
	if len(args) == 0 {
		cmd.name = cmd.Settings.ActiveProfile

		if len(cmd.name) == 0 {
			return errors.New("no profile currently active")
		}
	} else {
		cmd.name = strings.TrimSpace(args[0])

		if len(cmd.name) == 0 {
			return errors.New("profile name not specified")
		}
	}

	return nil
}

// Run executes the command
func (cmd *ShowCommand) Run() error {
	if len(cmd.Settings.Profiles) == 0 {
		return errors.New("no profiles currently exist")
	}

	profile, ok := cmd.Settings.Profiles[cmd.name]
	if !ok {
		return fmt.Errorf("profile '%s' was not found", cmd.name)
	}

	fmt.Fprintf(cmd.Out, "Name: %s\n", profile.Name)

	return nil
}
