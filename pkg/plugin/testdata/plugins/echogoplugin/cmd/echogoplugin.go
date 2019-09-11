/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

const (
	usage = "echogoplugin --message <message>"

	shortDesc = "Echos back the message provided by the --message flag"

	longDesc = `
The echogoplugin command demonstrates the addition of new commands to fabric-cli using Go plugins. Commands defined
as Go plugins can contain custom flags, long descriptions and examples.

The echogoplugin command is a simple command that echos back the message provided by the --message flag.`

	example = `
> fabric help echogoplugin

... prints out a long description of the echogoplugin command along with example usages

> fabric echogoplugin --message "Hello World!"

... prints out:

Hello World!
`
)

// newCmd returns the command
func newCmd(settings *environment.Settings) *cobra.Command {
	c := &echoCommand{}
	c.Settings = settings

	cmd := &cobra.Command{
		Use:     usage,
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return c.validate()
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}
	cmd.Flags().StringVar(&c.message, "message", "", "sets the message to echo")
	return cmd
}

// New returns a new command
func New(settings *environment.Settings) *cobra.Command {
	return newCmd(settings)
}

// echoCommand implements a Go plugin command
type echoCommand struct {
	common.Command

	message string
}

func (c *echoCommand) validate() error {
	if c.message == "" {
		return fmt.Errorf("no message specified")
	}
	return nil
}

func (c *echoCommand) run() {
	_, err := fmt.Fprintln(c.Settings.Streams.Out, c.message)
	if err != nil {
		panic(err.Error())
	}
}
