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

// NewProfileUseCommand creates a new "fabric profile use" command
func NewProfileUseCommand(settings *environment.Settings) *cobra.Command {
	pcmd := profileUseCommand{
		out: settings.Streams.Out,
	}

	cmd := &cobra.Command{
		Use:   "use <profilename>",
		Short: "change active profile",
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

type profileUseCommand struct {
	out    io.Writer
	config *environment.Settings

	name string
}

func (cmd *profileUseCommand) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("profile name not specified")
	}

	cmd.name = strings.TrimSpace(args[0])

	if len(cmd.name) == 0 {
		return errors.New("profile name not specified")
	}

	return nil
}

func (cmd *profileUseCommand) run() error {
	var found bool
	for _, p := range cmd.config.Profiles {
		if p.Name == cmd.name {
			found = true
			cmd.config.ActiveProfile = cmd.name
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' was not found", cmd.name)
	}

	err := cmd.config.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.out, "successfully set active profile to '%s'\n", cmd.name)

	return nil
}
