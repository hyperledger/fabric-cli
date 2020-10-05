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

// NewQueryCommittedCommand creates a new "fabric lifecycle querycommitted" command
func NewQueryCommittedCommand(settings *environment.Settings) *cobra.Command {
	c := QueryCommittedCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "querycommitted <chaincode-name>",
		Short: "Query for committed chaincodes",
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

	c.AddArg(&c.ChaincodeName)

	flags := cmd.Flags()
	flags.StringVar(&c.OutputFormat, "output", "", outputFormatUsage)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// QueryCommittedCommand implements the chaincode queryapproved command
type QueryCommittedCommand struct {
	BaseCommand

	ChaincodeName string
	OutputFormat  string
}

// Validate checks the required parameters for run
func (c *QueryCommittedCommand) Validate() error {
	if c.ChaincodeName == "" {
		return errors.New("chaincode name not specified")
	}

	return nil
}

// Run executes the command
func (c *QueryCommittedCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	committedChaincodes, err := c.ResourceManagement.LifecycleQueryCommittedCC(
		context.Channel,
		resmgmt.LifecycleQueryCommittedCCRequest{
			Name: c.ChaincodeName,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(context.Peers[0]),
	)
	if err != nil {
		return err
	}

	if c.OutputFormat == jsonFormat {
		return c.printJSONResponse(committedChaincodes)
	}

	c.printResponse(committedChaincodes)

	return nil
}

func (c *QueryCommittedCommand) printResponse(defs []resmgmt.LifecycleChaincodeDefinition) {
	if len(defs) == 0 {
		c.println("No committed chaincodes")

		return
	}

	for _, def := range defs {
		var approvingOrgs []string
		var nonApprovingOrgs []string

		for org, approved := range def.Approvals {
			if approved {
				approvingOrgs = append(approvingOrgs, org)
			} else {
				nonApprovingOrgs = append(nonApprovingOrgs, org)
			}
		}

		c.printf("Name: %s, Version: %s, Sequence: %d, Validation Plugin: %s, Endorsement Plugin: %s,"+
			" Channel Config Policy: %s, Init Required: %t, Approving orgs: %s, Non-approving orgs: %s\n",
			def.Name, def.Version, def.Sequence, def.ValidationPlugin, def.EndorsementPlugin, def.ChannelConfigPolicy,
			def.InitRequired, approvingOrgs, nonApprovingOrgs)

		for _, collConfig := range def.CollectionConfig {
			cfg := collConfig.GetStaticCollectionConfig()

			c.printf("- Collection: %s, Blocks to Live: %d, Maximum Peer Count: %d,"+
				" Required Peer Count: %d, MemberOnlyRead: %t, cfg.MemberOnlyWrite: %t\n",
				cfg.Name, cfg.BlockToLive, cfg.MaximumPeerCount, cfg.RequiredPeerCount, cfg.MemberOnlyRead, cfg.MemberOnlyWrite)
		}
	}
}
