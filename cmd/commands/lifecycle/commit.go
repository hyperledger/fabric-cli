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

// NewCommitCommand creates a new "fabric lifecycle commit" command
func NewCommitCommand(settings *environment.Settings) *cobra.Command {
	c := CommitCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "commit <chaincode-name> <version> <sequence>",
		Short: "commit a chaincode",
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
	c.AddArg(&c.Sequence)

	flags := cmd.Flags()
	flags.StringVar(&c.SignaturePolicy, "policy", "", "sets the signature policy")
	flags.StringVar(&c.ChannelConfigPolicy, "channel-config-policy", "", "sets the channel config policy")
	flags.StringVar(&c.CollectionsConfig, "collections-config", "", "sets the path to the collections config file")
	flags.BoolVar(&c.InitRequired, "init-required", false, "indicates whether the chaincode requires 'Init' to be invoked")
	flags.StringVar(&c.EndorsementPlugin, "endorsement-plugin", "", "sets the endorsement plugin")
	flags.StringVar(&c.ValidationPlugin, "validation-plugin", "", "sets the validation plugin")
	flags.StringArrayVar(&c.Peers, "peer", []string{}, "sets a peer to which to send the commit (this option may be specified multiple times)")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// CommitCommand implements the lifecycle commit command
type CommitCommand struct {
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
	Peers               []string
}

// Validate checks the required parameters for run
func (c *CommitCommand) Validate() error {
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
func (c *CommitCommand) Run() error {
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

	req := resmgmt.LifecycleCommitCCRequest{
		Name:                c.Name,
		Version:             c.Version,
		Sequence:            sequence,
		SignaturePolicy:     signaturePolicy,
		ChannelConfigPolicy: c.ChannelConfigPolicy,
		CollectionConfig:    collectionsConfig,
		InitRequired:        c.InitRequired,
		EndorsementPlugin:   c.EndorsementPlugin,
		ValidationPlugin:    c.ValidationPlugin,
	}

	peers := c.Peers
	if len(peers) == 0 {
		peers = context.Peers
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}

	if _, err := c.ResourceManagement.LifecycleCommitCC(context.Channel, req, options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully committed chaincode '%s'\n", c.Name)

	return nil
}
