/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChaincodeEventsCommand creates a new "fabric chaincode events" command
func NewChaincodeEventsCommand(settings *environment.Settings) *cobra.Command {
	c := EventsCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "events <chaincode-name>",
		Short: "listen for chaincode events",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := c.Complete(); err != nil {
				return err
			}

			if err := c.Validate(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.ChaincodeName)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// EventsCommand implements the chaincode events command
type EventsCommand struct {
	BaseCommand

	ChaincodeName string
}

// Validate checks the required parameters for run
func (c *EventsCommand) Validate() error {
	if len(c.ChaincodeName) == 0 {
		return errors.New("chaincode name not specified")
	}

	return nil
}

// Run executes the command
func (c *EventsCommand) Run() error {
	registration, eventCh, err := c.Channel.RegisterChaincodeEvent(c.ChaincodeName, "")
	if err != nil {
		return err
	}

	defer c.Channel.UnregisterChaincodeEvent(registration)

	for event := range eventCh {
		fmt.Fprintln(c.Settings.Streams.Out, string(event.Payload))
	}

	return nil
}
