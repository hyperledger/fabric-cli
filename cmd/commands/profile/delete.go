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

// NewProfileDeleteCommand creates a new "fabric profile delete" command
func NewProfileDeleteCommand(settings *environment.Settings) *cobra.Command {
	pcmd := profileDeleteCommand{
		out: settings.Streams.Out,
	}

	cmd := &cobra.Command{
		Use:   "delete <profilename>",
		Short: "delete a configuration profile",
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

type profileDeleteCommand struct {
	out    io.Writer
	config *environment.Settings

	name string
}

func (cmd *profileDeleteCommand) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("profile name not specified")
	}

	cmd.name = strings.TrimSpace(args[0])

	if len(cmd.name) == 0 {
		return errors.New("profile name not specified")
	}

	return nil
}

func (cmd *profileDeleteCommand) run() error {
	var found bool
	for i, p := range cmd.config.Profiles {
		if p.Name == cmd.name {
			found = true
			cmd.config.Profiles = append(cmd.config.Profiles[:i], cmd.config.Profiles[i+1:]...)
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' was not found", cmd.name)
	}

	if cmd.name == cmd.config.ActiveProfile {
		cmd.config.ActiveProfile = ""
	}

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.out, "successfully deleted profile '%s'\n", cmd.name)

	return nil
}
