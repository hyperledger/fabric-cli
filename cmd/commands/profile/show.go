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

// NewProfileShowCommand creates a new "fabric profile show" command
func NewProfileShowCommand(settings *environment.Settings) *cobra.Command {
	pcmd := profileShowCommand{
		out:      settings.Streams.Out,
		profiles: settings.Profiles,
		active:   settings.ActiveProfile,
	}

	cmd := &cobra.Command{
		Use:   "show [profilename]",
		Short: "show the metadata of the active profile",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type profileShowCommand struct {
	out      io.Writer
	profiles []*environment.Profile
	active   string

	name string
}

func (cmd *profileShowCommand) complete(args []string) error {
	if len(args) == 0 {
		cmd.name = cmd.active
	} else {
		cmd.name = strings.TrimSpace(args[0])
	}

	if len(cmd.name) == 0 {
		return errors.New("no profile currently active")
	}

	return nil
}

func (cmd *profileShowCommand) run() error {
	if len(cmd.profiles) == 0 {
		return errors.New("no profiles currently exist")
	}

	for _, p := range cmd.profiles {
		if p.Name == cmd.name {
			fmt.Fprintf(cmd.out, "Name: %s\n", p.Name)
		}
	}

	return nil
}
