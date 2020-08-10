/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewCommitCommand creates a new "fabric chaincode approve" command
func NewApproveCommand(settings *environment.Settings) *cobra.Command {
	c := ApproveCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "approve <chaincode-name> <version> <package-id> <sequence>",
		Short: "approve a chaincode for an org",
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

	c.AddArg(&c.Name)
	c.AddArg(&c.Version)
	c.AddArg(&c.PackageID)
	c.AddArg(&c.Sequence)

	flags := cmd.Flags()
	flags.StringVar(&c.SignaturePolicy, "policy", "", "sets the signature policy")
	flags.StringVar(&c.ChannelConfigPolicy, "channel-config-policy", "", "sets the channel config policy")
	flags.StringVar(&c.CollectionsConfig, "collections-config", "", "sets the path to the collections config file")
	flags.BoolVar(&c.InitRequired, "init-required", false, "indicates whether the chaincode requires 'Init' to be invoked")
	flags.StringVar(&c.EndorsementPlugin, "endorsement-plugin", "", "sets the endorsement plugin")
	flags.StringVar(&c.ValidationPlugin, "validation-plugin", "", "sets the validation plugin")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// ApproveCommand implements the chaincode instantiate command
type ApproveCommand struct {
	BaseCommand

	Name                string
	Version             string
	PackageID           string
	SignaturePolicy     string
	ChannelConfigPolicy string
	CollectionsConfig   string
	Sequence            string
	InitRequired        bool
	EndorsementPlugin   string
	ValidationPlugin    string
}

// Validate checks the required parameters for run
func (c *ApproveCommand) Validate() error {
	if c.Name == "" {
		return errors.New("chaincode name not specified")
	}

	if c.Version == "" {
		return errors.New("chaincode version not specified")
	}

	if c.PackageID == "" {
		return errors.New("chaincode package ID not specified")
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
func (c *ApproveCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
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

	req := resmgmt.LifecycleApproveCCRequest{
		Name:                c.Name,
		Version:             c.Version,
		PackageID:           c.PackageID,
		Sequence:            sequence,
		SignaturePolicy:     signaturePolicy,
		ChannelConfigPolicy: c.ChannelConfigPolicy,
		CollectionConfig:    collectionsConfig,
		InitRequired:        c.InitRequired,
		EndorsementPlugin:   c.EndorsementPlugin,
		ValidationPlugin:    c.ValidationPlugin,
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}

	if _, err := c.ResourceManagement.LifecycleApproveCC(context.Channel, req, options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully approved chaincode '%s'\n", c.Name)

	return nil
}
