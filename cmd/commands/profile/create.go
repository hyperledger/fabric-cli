/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

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
	pcmd := profileCreateCommand{
		out: settings.Streams.Out,
	}

	cmd := &cobra.Command{
		Use:   "create <profilename>",
		Short: "create a new configuration profile",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config, err := settings.FromFile()
			if err != nil {
				return err
			}

			pcmd.config = config

			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type profileCreateCommand struct {
	out    io.Writer
	config *environment.Settings

	name string
}

func (cmd *profileCreateCommand) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("profile name not specified")
	}

	cmd.name = strings.TrimSpace(args[0])

	if len(cmd.name) == 0 {
		return errors.New("profile name not specified")
	}

	for _, p := range cmd.config.Profiles {
		if cmd.name == p.Name {
			return fmt.Errorf("profile '%s' already exists", cmd.name)
		}
	}

	return nil
}

func (cmd *profileCreateCommand) run() error {
	profile := &environment.Profile{
		Name: cmd.name,
	}

	cmd.config.Profiles = append(cmd.config.Profiles, profile)

	if len(cmd.config.ActiveProfile) == 0 {
		cmd.config.ActiveProfile = profile.Name
	}

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.out, "successfully created profile '%s'\n", cmd.name)

	return nil
}
