package chaincode

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/spf13/cobra"
)

// NewChaincodePackageCommand creates a new "fabric chaincode package" command
func NewChaincodePackageCommand(settings *environment.Settings) *cobra.Command {
	c := PackageCommand{}

	c.Settings = settings

	cmd := &cobra.Command{
		Use:   "package <chaincode-name> <path>",
		Short: "package a chaincode (only golang supported)",
		Args:  c.ParseArgs(),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return c.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Run()
		},
	}

	c.AddArg(&c.ChaincodeName)
	c.AddArg(&c.ChaincodePath)

	cmd.SetOutput(c.Settings.Streams.Out)

	return cmd
}

// PackageCommand implements the chaincode package command
type PackageCommand struct {
	BaseCommand

	ChaincodeName string
	ChaincodePath string
}

// Validate checks the required parameters for run
func (c *PackageCommand) Validate() error {
	if len(c.ChaincodeName) == 0 {
		return errors.New("chaincode name not specified")
	}

	if len(c.ChaincodePath) == 0 {
		return errors.New("chaincode path not specified")
	}

	return nil
}

// Run executes the command
func (c *PackageCommand) Run() error {
	pkg, err := gopackager.NewCCPackage(c.ChaincodePath, "")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./%s.tgz", c.ChaincodeName), pkg.Code, 0644); err != nil {
		return err
	}

	fmt.Fprintf(c.Settings.Streams.Out, "successfully packaged chaincode '%s'\n", c.ChaincodeName)

	return nil
}
