/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/common"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
)

const (
	jsonFormat = "json"

	outputFormatUsage = `The output format for query results. If set to 'json' then the response is output in JSON format,
otherwise the response is output in human-readable text.`
)

// NewCommand creates a new "fabric lifecycle" command
func NewCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Manage chaincode lifecycle",
	}

	cmd.AddCommand(
		NewPackageCommand(settings),
		NewInstallCommand(settings),
		NewApproveCommand(settings),
		NewCommitCommand(settings),
		NewQueryInstalledCommand(settings),
		NewQueryApprovedCommand(settings),
		NewCheckCommitReadinessCommand(settings),
		NewQueryCommittedCommand(settings),
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

func (c *BaseCommand) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(c.Settings.Streams.Out, format, a...)
	if err != nil {
		panic(err)
	}
}

func (c *BaseCommand) println(a ...interface{}) {
	_, err := fmt.Fprintln(c.Settings.Streams.Out, a...)
	if err != nil {
		panic(err)
	}
}

func (c *BaseCommand) print(a ...interface{}) {
	_, err := fmt.Fprint(c.Settings.Streams.Out, a...)
	if err != nil {
		panic(err)
	}
}

func (c *BaseCommand) printJSONResponse(v interface{}) error {
	respBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.print(string(respBytes))

	return nil
}
