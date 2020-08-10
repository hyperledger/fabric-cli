/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	lifecyclepkg "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// NewPackageCommand creates a new "fabric lifecycle chaincode package" command
func NewPackageCommand(settings *environment.Settings) *cobra.Command {
	c := PackageCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "package <chaincode-label> <chaincode-type> <path>",
		Short: "package a chaincode",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.Label)
	c.AddArg(&c.Type)
	c.AddArg(&c.Path)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// PackageCommand implements the chaincode package command
type PackageCommand struct {
	BaseCommand

	Path  string
	Label string
	Type  string
}

// Validate checks the required parameters for run
func (c *PackageCommand) Validate() error {
	if c.Label == "" {
		return errors.New("chaincode label not specified")
	}

	if c.Path == "" {
		return errors.New("chaincode path not specified")
	}

	if c.Type == "" {
		return errors.New("chaincode type not specified")
	}

	ccType, ok := pb.ChaincodeSpec_Type_value[strings.ToUpper(c.Type)]
	if !ok || ccType == int32(pb.ChaincodeSpec_UNDEFINED) {
		return errors.New("unsupported chaincode type")
	}

	return nil
}

// Run executes the command
func (c *PackageCommand) Run() error {
	pkgBytes, err := lifecyclepkg.NewCCPackage(&lifecyclepkg.Descriptor{
		Path:  c.Path,
		Type:  pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[strings.ToUpper(c.Type)]),
		Label: c.Label,
	})
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./%s.tgz", c.Label), pkgBytes, 0644); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully packaged chaincode '%s'\n", c.Label)

	return nil
}
