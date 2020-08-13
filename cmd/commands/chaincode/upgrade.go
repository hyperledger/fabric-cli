/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChaincodeUpgradeCommand creates a new "fabric chaincode upgrade" command
func NewChaincodeUpgradeCommand(settings *environment.Settings) *cobra.Command {
	c := UpgradeCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "upgrade <chaincode-name> <version> <path>",
		Short: "Upgrade a chaincode",
		Long:  "Upgrade a chaincode with chaincode-name version path args collection and policy",
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
	c.AddArg(&c.ChaincodeVersion)
	c.AddArg(&c.ChaincodePath)

	flags := cmd.Flags()
	flags.StringArrayVar(&c.ChaincodeArgs, "args", []string{}, "Set the upgrade arguments")
	flags.StringVar(&c.ChaincodePolicy, "policy", "", "Set the endorsement policy")
	flags.StringVar(&c.ChaincodeCollectionsConfig, "collections-config", "", "Set the path to the collections config file")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// UpgradeCommand implements the chaincode upgrade command
type UpgradeCommand struct {
	BaseCommand

	ChaincodeName    string
	ChaincodeVersion string
	ChaincodePath    string

	ChaincodeArgs              []string
	ChaincodePolicy            string
	ChaincodeCollectionsConfig string
}

// Validate checks the required parameters for run
func (c *UpgradeCommand) Validate() error {
	if len(c.ChaincodeName) == 0 {
		return errors.New("chaincode name not specified")
	}

	if len(c.ChaincodeVersion) == 0 {
		return errors.New("chaincode version not specified")
	}

	if len(c.ChaincodePath) == 0 {
		return errors.New("chaincode path not specified")
	}

	return nil
}

// Run executes the command
func (c *UpgradeCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	args, err := json.Marshal(c.ChaincodeArgs)
	if err != nil {
		return err
	}

	policy, err := common.GetChaincodePolicy(c.ChaincodePolicy)
	if err != nil {
		return err
	}

	collectionsConfig, err := common.GetCollectionConfigFromFile(c.ChaincodeCollectionsConfig)
	if err != nil {
		return err
	}

	req := resmgmt.UpgradeCCRequest{
		Name:       c.ChaincodeName,
		Path:       c.ChaincodePath,
		Version:    c.ChaincodeVersion,
		Args:       [][]byte{[]byte("init"), args},
		Policy:     policy,
		CollConfig: collectionsConfig,
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}

	if _, err := c.ResourceManagement.UpgradeCC(context.Channel, req, options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully upgraded chaincode '%s'\n", c.ChaincodeName)

	return nil
}
