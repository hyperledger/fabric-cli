/*
Copyright © 2019 TopJohn <xzj19922010@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package version

import (
	"fmt"
	"github.com/hyperledger/fabric-cli/common/metadata"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
	"runtime"
)

//Program name
const ProgramName = "fabric"

//NewVersionCommand creates a new "fabric version" command
func NewVersionCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print fabric cmd version.",
		Long:  `Print current version of fabric command line tool.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("trailing args detected")
			}
			cmd.SilenceUsage = true
			fmt.Print(GetMetaInfo())
			return nil
		},
	}
	cmd.SetOutput(settings.Streams.Out)
	return cmd
}

//GetMetaInfo returns version information for the fabric cmd.
func GetMetaInfo() string {
	return fmt.Sprintf("%s:\n Version: %s\n Commit SHA: %s\n Go Version: %s\n"+
		" OS/Arch: %s\n",
		ProgramName, metadata.Version, metadata.CommitSHA, runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
}
