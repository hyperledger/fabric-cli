/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

type fileWriter func(filename string, data []byte, perm os.FileMode) error

// NewGetInstalledPkgCommand creates a new "fabric lifecycle getinstalledpkg" command
func NewGetInstalledPkgCommand(settings *environment.Settings) *cobra.Command {
	c := GetInstalledPkgCommand{}

	c.Settings = settings
	c.WriteFile = ioutil.WriteFile

	cmd := &cobra.Command{
		Use:   "getinstalledpackage <peer> <package ID>",
		Short: "Get an installed chaincode package",
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
	c.AddArg(&c.PackageID)

	flags := cmd.Flags()
	flags.StringVar(&c.OutputDirectory, "output-directory", "",
		"sets the output directory for the chaincode package file (default is current directory)")

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// GetInstalledPkgCommand implements the chaincode getinstalledpackage command
type GetInstalledPkgCommand struct {
	BaseCommand

	Peer            string
	PackageID       string
	OutputDirectory string
	WriteFile       fileWriter
}

// Validate checks the required parameters for run
func (c *GetInstalledPkgCommand) Validate() error {
	if c.Peer == "" {
		return errors.New("peer not specified")
	}

	if c.PackageID == "" {
		return errors.New("package ID not specified")
	}

	return nil
}

// Run executes the command
func (c *GetInstalledPkgCommand) Run() error {
	pkgBytes, err := c.ResourceManagement.LifecycleGetInstalledCCPackage(
		c.PackageID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(c.Peer),
	)
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(filepath.Join(c.OutputDirectory, c.PackageID+".tar.gz"))
	if err != nil {
		return err
	}

	err = c.WriteFile(filePath, pkgBytes, os.ModePerm)
	if err != nil {
		return errors.WithMessagef(err, "failed to write chaincode package to file %s", filePath)
	}

	c.printf("Chaincode package saved to %s\n", filePath)

	return nil
}
