/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewHomeCommand is a fabric plugin for "fabric home"
func NewHomeCommand() *cobra.Command {
	hcmd := &homeCommand{
		out: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "home",
		Short: "output the current fabric home directory",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			settings, err := environment.GetSettings()
			if err != nil {
				return err
			}

			hcmd.home = settings.Home
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			hcmd.run()
		},
	}

	return cmd
}

type homeCommand struct {
	home environment.Home
	out  io.Writer
}

func (h *homeCommand) run() {
	fmt.Fprintln(h.out, h.home)
}

func main() {
	cmd := NewHomeCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
