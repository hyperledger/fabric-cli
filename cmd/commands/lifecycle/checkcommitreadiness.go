/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewCheckCommitReadinessCommand creates a new "fabric lifecycle checkcommitreadiness" command
func NewCheckCommitReadinessCommand(settings *environment.Settings) *cobra.Command {
	c := CheckCommitReadinessCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "checkcommitreadiness <chaincode-name> <version> <sequence>",
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

	c.AddArg(&c.Name)
	c.AddArg(&c.Version)
	c.AddArg(&c.Sequence)

	flags := cmd.Flags()
	flags.StringVar(&c.SignaturePolicy, "policy", "", "sets the signature policy")
	flags.StringVar(&c.ChannelConfigPolicy, "channel-config-policy", "", "sets the channel config policy")
	flags.StringVar(&c.CollectionsConfig, "collections-config", "", "sets the path to the collections config file")
	flags.BoolVar(&c.InitRequired, "init-required", false, "indicates whether the chaincode requires 'Init' to be invoked")
	flags.StringVar(&c.EndorsementPlugin, "endorsement-plugin", "", "sets the endorsement plugin")
	flags.StringVar(&c.ValidationPlugin, "validation-plugin", "", "sets the validation plugin")
	flags.StringVar(&c.OutputFormat, "output", "", outputFormatUsage)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// CheckCommitReadinessCommand implements the chaincode checkcommitreadiness command
type CheckCommitReadinessCommand struct {
	BaseCommand

	Name                string
	Version             string
	Sequence            string
	SignaturePolicy     string
	ChannelConfigPolicy string
	CollectionsConfig   string
	InitRequired        bool
	EndorsementPlugin   string
	ValidationPlugin    string
	OutputFormat        string
}

// Validate checks the required parameters for run
func (c *CheckCommitReadinessCommand) Validate() error {
	if c.Name == "" {
		return errors.New("chaincode name not specified")
	}

	if c.Version == "" {
		return errors.New("chaincode version not specified")
	}

	if c.Sequence == "" {
		return errors.New("sequence not specified")
	}

	sequence, err := strconv.ParseInt(c.Sequence, 10, 64)
	if err != nil {
		return errors.WithMessage(err, "invalid sequence")
	}

	if sequence <= 0 {
		return errors.New("sequence must be greater than 0")
	}

	return nil
}

// Run executes the command
func (c *CheckCommitReadinessCommand) Run() error {
	/*context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	signaturePolicy, err := common.GetChaincodePolicy(c.SignaturePolicy)
	if err != nil {
		return err
	}

	collectionsConfig, err := common.GetCollectionConfigFromFile(c.CollectionsConfig)
	if err != nil {
		return err
	}

	sequence, err := strconv.ParseInt(c.Sequence, 10, 64)
	if err != nil {
		return errors.WithMessage(err, "invalid sequence")
	}

	/*resp, err := c.ResourceManagement.LifecycleCheckCCCommitReadiness(
		context.Channel,
		resmgmt.LifecycleCheckCCCommitReadinessRequest{
			Name:                c.Name,
			Version:             c.Version,
			Sequence:            sequence,
			EndorsementPlugin:   c.EndorsementPlugin,
			ValidationPlugin:    c.ValidationPlugin,
			SignaturePolicy:     signaturePolicy,
			ChannelConfigPolicy: c.ChannelConfigPolicy,
			CollectionConfig:    collectionsConfig,
			InitRequired:        c.InitRequired,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(context.Peers[0]),
	)
	if err != nil {
		return err
	}

	if c.OutputFormat == jsonFormat {
		return c.printJSONResponse(resp)
	}

	c.printResponse(resp)*/

	return nil
}

/*
func (c *CheckCommitReadinessCommand) printResponse(crr resmgmt.LifecycleCheckCCCommitReadinessResponse) {
	var approvingOrgs []string
	var nonApprovingOrgs []string

	for org, approved := range crr.Approvals {
		if approved {
			approvingOrgs = append(approvingOrgs, org)
		} else {
			nonApprovingOrgs = append(nonApprovingOrgs, org)
		}
	}

	c.printf("Approving orgs: %s\n", approvingOrgs)
	c.printf("Non-approving orgs: %s\n", nonApprovingOrgs)
}
*/
