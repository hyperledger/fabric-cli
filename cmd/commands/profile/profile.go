/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package profile

import (
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
)

// NewProfileCommand creates a new "fabric profile" command
func NewProfileCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}

	cmd.AddCommand(
		NewProfileListCommand(settings),
		NewProfileUseCommand(settings),
		NewProfileShowCommand(settings),
		NewProfileCreateCommand(settings),
		NewProfileDeleteCommand(settings),
	)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}
