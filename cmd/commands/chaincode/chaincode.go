/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

// NewChaincodeCommand creates a new "fabric chaincode" command
func NewChaincodeCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chaincode",
		Short: "Manage chaincode",
	}

	cmd.AddCommand(
		NewChaincodeListCommand(settings),
		NewChaincodePackageCommand(settings),
		NewChaincodeInstallCommand(settings),
		NewChaincodeInstantiateCommand(settings),
		NewChaincodeUpgradeCommand(settings),
		NewChaincodeQueryCommand(settings),
		NewChaincodeInvokeCommand(settings),
		NewChaincodeEventsCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}

// BaseCommand implements common channel command functions
type BaseCommand struct {
	common.Command

	Factory            fabric.Factory
	Channel            fabric.Channel
	ResourceManagement fabric.ResourceManagement
}

// Complete initializes all clients needed for Run
func (c *BaseCommand) Complete() error {
	var err error

	if c.Factory == nil {
		c.Factory, err = fabric.NewFactory(c.Settings.Config)
		if err != nil {
			return err
		}
	}

	c.Channel, err = c.Factory.Channel()
	if err != nil {
		return err
	}

	c.ResourceManagement, err = c.Factory.ResourceManagement()
	if err != nil {
		return err
	}

	return nil
}
