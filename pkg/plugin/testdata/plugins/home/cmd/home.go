/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewHomeCommand is a fabric plugin for "fabric home"
func NewHomeCommand() *cobra.Command {
	c := &HomeCommand{}

	cmd := &cobra.Command{
		Use:   "home",
		Short: "output the current fabric home directory",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			settings := environment.NewDefaultSettings()

			if err := settings.Init(&pflag.FlagSet{}); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}

			c.Settings = settings

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			c.run()
		},
	}

	return cmd
}

// HomeCommand implements a home command
type HomeCommand struct {
	common.Command
}

func (c *HomeCommand) run() {
	fmt.Fprintln(c.Settings.Streams.Out, c.Settings.Home)
}

func main() {
	cmd := NewHomeCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
