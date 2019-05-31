/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChaincodeListCommand creates a new "fabric chaincode list" command
func NewChaincodeListCommand(settings *environment.Settings) *cobra.Command {
	c := ListCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all chaincodes",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Complete()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&c.Installed, "installed", false, "include chaincode installed on peer's filesystem")
	flags.BoolVar(&c.Instantiated, "instantiated", false, "include instantiated chaincode")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ListCommand implements the chaincode list command
type ListCommand struct {
	BaseCommand

	Installed    bool
	Instantiated bool
}

// Run executes the command
func (c *ListCommand) Run() error {
	if !c.Installed && !c.Instantiated {
		c.Installed = true
		c.Instantiated = true
	}

	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
	}

	if c.Installed {
		resp, err := c.ResourceManagement.QueryInstalledChaincodes(options...)
		if err != nil {
			return err
		}

		fmt.Fprintln(c.Settings.Streams.Out, "Installed Chaincode:")
		for _, chaincode := range resp.Chaincodes {
			fmt.Fprintf(c.Settings.Streams.Out, " - %s\n", chaincode.Name)
		}
	}

	if c.Instantiated {
		resp, err := c.ResourceManagement.QueryInstalledChaincodes(options...)
		if err != nil {
			return err
		}

		fmt.Fprintln(c.Settings.Streams.Out, "Instantiated Chaincode:")
		for _, chaincode := range resp.Chaincodes {
			fmt.Fprintf(c.Settings.Streams.Out, " - %s\n", chaincode.Name)
		}
	}

	return nil
}
