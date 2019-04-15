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

// NewProfileDeleteCommand creates a new "fabric profile delete" command
func NewProfileDeleteCommand(settings *environment.Settings) *cobra.Command {
	c := DeleteCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "delete <profile-name>",
		Short: "delete a Configuration profile",
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

// DeleteCommand implements the profile delete command
type DeleteCommand struct {
	Out      io.Writer
	Settings *environment.Settings

	name   string
	config *environment.Settings
}

// Complete populates required fields for Run
func (cmd *DeleteCommand) Complete(args []string) error {
	config, err := cmd.Settings.FromFile()
	if err != nil {
		return err
	}

	cmd.config = config

	if len(args) == 0 {
		return errors.New("profile name not specified")
	}

	cmd.name = strings.TrimSpace(args[0])

	if len(cmd.name) == 0 {
		return errors.New("profile name not specified")
	}

	return nil
}

// Run executes the command
func (cmd *DeleteCommand) Run() error {
	if _, ok := cmd.config.Profiles[cmd.name]; !ok {
		return fmt.Errorf("profile '%s' was not found", cmd.name)
	}

	delete(cmd.config.Profiles, cmd.name)

	if cmd.name == cmd.config.ActiveProfile {
		cmd.config.ActiveProfile = ""
	}

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Out, "successfully deleted profile '%s'\n", cmd.name)

	return nil
}
