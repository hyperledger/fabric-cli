/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewInstallCommand creates a new "fabric lifecycle install" command
func NewInstallCommand(settings *environment.Settings) *cobra.Command {
	c := InstallCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "install <chaincode-label> <path>",
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

	c.AddArg(&c.Label)
	c.AddArg(&c.Path)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// InstallCommand implements the chaincode install command
type InstallCommand struct {
	BaseCommand

	Label string
	Path  string
}

// Validate checks the required parameters for run
func (c *InstallCommand) Validate() error {
	if c.Label == "" {
		return errors.New("chaincode label not specified")
	}

	if c.Path == "" {
		return errors.New("chaincode path not specified")
	}

	return nil
}

// Run executes the command
func (c *InstallCommand) Run() error {
	/*context, err := c.Settings.Config.GetCurrentContext()
	if err != nil {
		return err
	}

	pkg, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return err
	}

	responses, err := c.ResourceManagement.LifecycleInstallCC(
		resmgmt.LifecycleInstallCCRequest{
			Label:   c.Label,
			Package: pkg,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(context.Peers...),
	)
	if err != nil {
		return err
	}

	if len(responses) == 0 {
		packageID := lifecycle.ComputePackageID(c.Label, pkg)
		fmt.Fprintf(c.Settings.Streams.Out, "chaincode '%s' has already been installed on all peers. Package ID '%s'\n", c.Label, packageID)
	} else {
		fmt.Fprintf(c.Settings.Streams.Out, "successfully installed chaincode '%s'. Package ID '%s'\n", c.Label, responses[0].PackageID)
	}
	*/
	return nil
}
