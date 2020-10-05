/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewQueryApprovedCommand creates a new "fabric lifecycle queryapproved" command
func NewQueryApprovedCommand(settings *environment.Settings) *cobra.Command {
	c := QueryApprovedCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "queryapproved <chaincode-name> <sequence>",
		Short: "Query for approved chaincodes",
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
	c.AddArg(&c.Sequence)

	flags := cmd.Flags()
	flags.StringVar(&c.OutputFormat, "output", "", outputFormatUsage)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// QueryApprovedCommand implements the chaincode queryapproved command
type QueryApprovedCommand struct {
	BaseCommand

	ChaincodeName string
	Sequence      string
	OutputFormat  string
}

// Validate checks the required parameters for run
func (c *QueryApprovedCommand) Validate() error {
	if c.ChaincodeName == "" {
		return errors.New("chaincode name not specified")
	}

	_, err := strconv.ParseInt(c.Sequence, 10, 64)
	if err != nil {
		return errors.WithMessage(err, "invalid sequence")
	}

	return nil
}

// Run executes the command
func (c *QueryApprovedCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	sequence, err := strconv.ParseInt(c.Sequence, 10, 64)
	if err != nil {
		return errors.WithMessage(err, "invalid sequence")
	}

	approvedChaincode, err := c.ResourceManagement.LifecycleQueryApprovedCC(
		context.Channel,
		resmgmt.LifecycleQueryApprovedCCRequest{
			Name:     c.ChaincodeName,
			Sequence: sequence,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(context.Peers[0]),
	)
	if err != nil {
		return err
	}

	if c.OutputFormat == jsonFormat {
		return c.printJSONResponse(approvedChaincode)
	}

	c.printResponse(approvedChaincode)

	return nil
}

func (c *QueryApprovedCommand) printResponse(ac resmgmt.LifecycleApprovedChaincodeDefinition) {
	c.printf("Name: %s, Version: %s, Package ID: %s, Sequence: %d, Validation Plugin: %s,"+
		" Endorsement Plugin: %s, Channel Config Policy: %s, Init Required: %t\n",
		ac.Name, ac.Version, ac.PackageID, ac.Sequence, ac.ValidationPlugin, ac.EndorsementPlugin, ac.ChannelConfigPolicy, ac.InitRequired)

	for _, collConfig := range ac.CollectionConfig {
		cfg := collConfig.GetStaticCollectionConfig()

		c.printf("- Collection: %s, Blocks to Live: %d, Maximum Peer Count: %d,"+
			" Required Peer Count: %d, MemberOnlyRead: %t, cfg.MemberOnlyWrite: %t\n",
			cfg.Name, cfg.BlockToLive, cfg.MaximumPeerCount, cfg.RequiredPeerCount, cfg.MemberOnlyRead, cfg.MemberOnlyWrite)
	}
}
