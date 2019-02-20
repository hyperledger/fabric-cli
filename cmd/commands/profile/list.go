/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

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
	pcmd := profileListCommand{
		out:      settings.Streams.Out,
		profiles: settings.Profiles,
		active:   settings.ActiveProfile,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all configuration profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}

	return cmd
}

type profileListCommand struct {
	out      io.Writer
	profiles []*environment.Profile
	active   string
}

func (cmd *profileListCommand) run() error {
	if len(cmd.profiles) == 0 {
		return errors.New("no profiles currently exist")
	}

	for _, p := range cmd.profiles {
		fmt.Fprint(cmd.out, p.Name)
		if p.Name == cmd.active {
			fmt.Fprint(cmd.out, " (active)")
		}
		fmt.Fprintln(cmd.out)

	}

	return nil
}
