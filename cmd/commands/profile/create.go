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

// NewProfileCreateCommand creates a new "fabric profile create" command
func NewProfileCreateCommand(settings *environment.Settings) *cobra.Command {
	c := CreateCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "create <profile-name>",
		Short: "create a new configuration profile",
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

// CreateCommand implements the profile create command
type CreateCommand struct {
	Out      io.Writer
	Settings *environment.Settings

	name   string
	config *environment.Settings
}

// Complete populates required fields for Run
func (cmd *CreateCommand) Complete(args []string) error {
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
func (cmd *CreateCommand) Run() error {
	if _, ok := cmd.config.Profiles[cmd.name]; ok {
		return fmt.Errorf("profile '%s' already exists", cmd.name)
	}

	cmd.config.Profiles[cmd.name] = &environment.Profile{
		Name: cmd.name,
	}

	if len(cmd.config.ActiveProfile) == 0 {
		cmd.config.ActiveProfile = cmd.name
	}

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Out, "successfully created profile '%s'\n", cmd.name)

	return nil
}
