/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// NewChaincodeInvokeCommand creates a new "fabric chaincode invoke" command
func NewChaincodeInvokeCommand(settings *environment.Settings) *cobra.Command {
	c := InvokeCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "invoke <chaincode-name>",
		Short: "invoke a chaincode",
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

	flags := cmd.Flags()
	flags.StringVar(&c.ChaincodeFcn, "fcn", "", "set the invoke function")
	flags.StringArrayVar(&c.ChaincodeArgs, "args", []string{}, "set the invoke arguments")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// InvokeCommand implements the chaincode invoke command
type InvokeCommand struct {
	BaseCommand

	ChaincodeName string

	ChaincodeFcn  string
	ChaincodeArgs []string
}

// Validate checks the required parameters for run
func (c *InvokeCommand) Validate() error {
	if len(c.ChaincodeName) == 0 {
		return errors.New("chaincode name not specified")
	}

	return nil
}

// Run executes the command
func (c *InvokeCommand) Run() error {
	args, err := json.Marshal(c.ChaincodeArgs)
	if err != nil {
		return err
	}

	req := channel.Request{
		ChaincodeID: c.ChaincodeName,
		Fcn:         c.ChaincodeFcn,
		Args:        [][]byte{args},
	}

	resp, err := c.Channel.Execute(req)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, string(resp.Payload))

	return nil
}
