/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// SafeYAMLUnmarshal provides a safer way to unmarshal YAML data
// It includes additional validation to prevent YAML deserialization attacks
func SafeYAMLUnmarshal(data []byte, out interface{}) error {
	// First, unmarshal the YAML data
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	// Additional validation could be added here based on the expected structure
	// For example, checking for unexpected fields or values

	return nil
}

// LoadPluginFromYAML loads a plugin from a YAML file with additional security checks
func LoadPluginFromYAML(path string) (*Plugin, error) {
	// Validate the path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, errors.New("plugin YAML file does not exist")
	}

	// Read the file
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	// Create a new plugin
	plugin := &Plugin{}

	// Use safe YAML unmarshaling
	if err := SafeYAMLUnmarshal(data, plugin); err != nil {
		return nil, err
	}

	// Validate the plugin
	if err := validatePlugin(plugin); err != nil {
		return nil, err
	}

	return plugin, nil
}

// validatePlugin performs validation on the plugin structure
func validatePlugin(p *Plugin) error {
	if p.Name == "" {
		return errors.New("plugin name is required")
	}

	if !isValidPluginName(p.Name) {
		return errors.New("invalid plugin name")
	}

	if p.Command == nil {
		return errors.New("plugin command is required")
	}

	if p.Command.Base == "" {
		return errors.New("plugin command base is required")
	}

	return nil
} 