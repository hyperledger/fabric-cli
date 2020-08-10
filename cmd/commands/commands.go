/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands/chaincode"
	"github.com/hyperledger/fabric-cli/cmd/commands/channel"
	"github.com/hyperledger/fabric-cli/cmd/commands/context"
	"github.com/hyperledger/fabric-cli/cmd/commands/lifecycle"
	"github.com/hyperledger/fabric-cli/cmd/commands/network"
	"github.com/hyperledger/fabric-cli/cmd/commands/plugin"
	"github.com/hyperledger/fabric-cli/cmd/commands/version"
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// All returns all subcommands that will be added to the root command
// Settings can be leveraged here to disable commands
func All(settings *environment.Settings) []*cobra.Command {
	return []*cobra.Command{
		// fabric plugin [subcommand]
		plugin.NewPluginCommand(settings),

		// fabric network [subcommand]
		network.NewNetworkCommand(settings),

		// fabric context [subcommand]
		context.NewContextCommand(settings),

		// fabric channel [subcommand]
		channel.NewChannelCommand(settings),

		// fabric channel [subcommand]
		chaincode.NewChaincodeCommand(settings),

		// fabric version
		version.NewVersionCommand(settings),

		// fabric lifecycle [subcommand]
		lifecycle.NewCommand(settings),
	}
}
