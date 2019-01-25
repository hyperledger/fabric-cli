/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/hyperledger/fabric-cli/cmd/commands/plugin"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// All returns all subcommands that will be added to the root command
// Settings can be leveraged here to disable commands
func All(settings *environment.Settings) []*cobra.Command {
	return []*cobra.Command{
		// fabric plugin < list | install | uninstall >
		plugin.NewPluginCommand(settings),
	}
}
