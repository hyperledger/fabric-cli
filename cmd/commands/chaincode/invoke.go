/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChaincodeInvokeCommand creates a new "fabric chaincode invoke" command
func NewChaincodeInvokeCommand(settings *environment.Settings) *cobra.Command {
	c := InvokeCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "invoke <chaincode-name>",
		Short: "Invoke a chaincode",
		Long:  "Invoke a chaincode with chaincode-name args function",
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
	flags.BoolVar(&c.IsInit, "is-init", false, "indicates whether or not this invocation is meant to initialize the chaincode")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// InvokeCommand implements the chaincode invoke command
type InvokeCommand struct {
	BaseCommand

	ChaincodeName string

	ChaincodeFcn  string
	ChaincodeArgs []string
	IsInit        bool
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
	/*fcn := c.ChaincodeFcn
	if c.IsInit {
		fcn = "Init"
	}

	req := channel.Request{
		ChaincodeID: c.ChaincodeName,
		Fcn:         fcn,
		Args:        common.AsByteArgs(c.ChaincodeArgs),
		IsInit:      c.IsInit,
	}

	resp, err := c.Channel.Execute(req, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Settings.Streams.Out, string(resp.Payload))
	*/
	return nil
}
