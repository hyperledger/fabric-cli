/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"strings"
)

// DefaultFilename is the filename for the plugin metadata
const DefaultFilename = "plugin.yaml"

// Plugin is an installed, third-party command
// Consider using "github.com/mitchellh/mapstructure" when the yaml structure
// gets more complex
type Plugin struct {
	Name        string   `yaml:"name"`
	Usage       string   `yaml:"usage"`
	Description string   `yaml:"description"`
	Command     *Command `yaml:"command"`

	Path string `yaml:"-"`
}

// Command is a parsed plugin command
type Command struct {
	Base string
	Args []string
}

// UnmarshalYAML implements the yaml unmarshaller interface
// This function deconstructs the command to base command and args
func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cmd string
	err := unmarshal(&cmd)
	if err != nil {
		return err
	}

	parts := strings.Split(cmd, " ")

	if len(parts) > 0 {
		c.Base = parts[0]
		c.Args = parts[1:]
	}

	return nil
}
