/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewQueryInstalledCommand creates a new "fabric lifecycle queryinstalled" command
func NewQueryInstalledCommand(settings *environment.Settings) *cobra.Command {
	c := QueryInstalledCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "queryinstalled <peer>",
		Short: "Query a peer for installed chaincodes",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if err := c.Complete(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Peer)

	flags := cmd.Flags()
	flags.StringVar(&c.OutputFormat, "output", "", outputFormatUsage)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// QueryInstalledCommand implements the chaincode queryinstalled command
type QueryInstalledCommand struct {
	BaseCommand

	Peer         string
	OutputFormat string
}

// Validate checks the required parameters for run
func (c *QueryInstalledCommand) Validate() error {
	if c.Peer == "" {
		return errors.New("peer not specified")
	}

	return nil
}

// Run executes the command
func (c *QueryInstalledCommand) Run() error {
	installedChaincodes, err := c.ResourceManagement.LifecycleQueryInstalledCC(
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(c.Peer),
	)
	if err != nil {
		return err
	}

	if c.OutputFormat == jsonFormat {
		return c.printJSONResponse(installedChaincodes)
	}

	return c.printResponse(installedChaincodes)
}

func (c *QueryInstalledCommand) printResponse(installedChaincodes []resmgmt.LifecycleInstalledCC) error {
	if len(installedChaincodes) == 0 {
		c.printf("No installed chaincodes on peer %s", c.Peer)

		return nil
	}

	c.println("Installed chaincodes:")

	for _, cc := range installedChaincodes {
		c.printf("- Package ID: %s, Label: %s\n", cc.PackageID, cc.Label)

		c.printReferences(cc.References)
	}

	return nil
}

func (c *QueryInstalledCommand) printReferences(refs map[string][]resmgmt.CCReference) {
	for channelID, refs := range refs {
		c.printf("-- References for channel [%s]:\n", channelID)

		for _, ref := range refs {
			c.printf("--- Name: %s, Version: %s\n", ref.Name, ref.Version)
		}
	}
}
