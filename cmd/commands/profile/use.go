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

// NewProfileUseCommand creates a new "fabric profile use" command
func NewProfileUseCommand(settings *environment.Settings) *cobra.Command {
	c := UseCommand{
		Out:      settings.Streams.Out,
		Settings: settings,
	}

	cmd := &cobra.Command{
		Use:   "use <profilename>",
		Short: "change active profile",
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

// UseCommand implements the profile show command
type UseCommand struct {
	Out      io.Writer
	Settings *environment.Settings

	config *environment.Settings
	name   string
}

// Complete populates required fields for Run
func (cmd *UseCommand) Complete(args []string) error {
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
func (cmd *UseCommand) Run() error {
	if _, ok := cmd.config.Profiles[cmd.name]; !ok {
		return fmt.Errorf("profile '%s' was not found", cmd.name)
	}

	cmd.config.ActiveProfile = cmd.name

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Out, "successfully set active profile to '%s'\n", cmd.name)

	return nil
}
