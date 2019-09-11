/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewChaincodeInstallCommand creates a new "fabric chaincode install" command
func NewChaincodeInstallCommand(settings *environment.Settings) *cobra.Command {
	c := InstallCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "install <chaincode-name> <version> <path>",
		Short: "install a chaincode",
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

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// InstallCommand implements the chaincode install command
type InstallCommand struct {
	BaseCommand

	ChaincodeName    string
	ChaincodeVersion string
	ChaincodePath    string
}

// Validate checks the required parameters for run
func (c *InstallCommand) Validate() error {
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
func (c *InstallCommand) Run() error {
	context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	pkg, err := ioutil.ReadFile(c.ChaincodePath)
	if err != nil {
		return err
	}

	req := resmgmt.InstallCCRequest{
		Name:    c.ChaincodeName,
		Path:    c.ChaincodePath,
		Version: c.ChaincodeVersion,
		Package: &resource.CCPackage{
			Type: peer.ChaincodeSpec_GOLANG,
			Code: pkg,
		},
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}

	if _, err := c.ResourceManagement.InstallCC(req, options...); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully installed chaincode '%s'\n", c.ChaincodeName)

	return nil
}
