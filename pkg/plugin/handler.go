/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	goplugin "plugin"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const pluginFactoryMethod = "New"

// ErrNotAGoPlugin is returned when attempting to load a file that's not a Go plugin
var ErrNotAGoPlugin = errors.New("not a Go plugin")

// ErrInvalidPluginPath is returned when a plugin path is invalid or potentially malicious
var ErrInvalidPluginPath = errors.New("invalid plugin path")

// Handler defines the required actions for managing plugins
type Handler interface {
	GetPlugins() ([]*Plugin, error)
	InstallPlugin(path string) error
	UninstallPlugin(name string) error
	LoadGoPlugin(path string, settings *environment.Settings) (*cobra.Command, error)
}

// DefaultHandler is the default plugin handler
type DefaultHandler struct {
	Dir      string
	Filename string
}

// GetPlugins returns all installed plugins found in the plugins directory
func (h *DefaultHandler) GetPlugins() ([]*Plugin, error) {
	// Ensure the plugins directory exists
	if _, err := os.Stat(h.Dir); os.IsNotExist(err) {
		return []*Plugin{}, nil
	}

	dirs, err := filepath.Glob(filepath.Join(h.Dir, "*"))
	if err != nil {
		return nil, err
	}

	plugins := []*Plugin{}
	if dirs == nil {
		return plugins, nil
	}

	for _, dir := range dirs {
		plugin, err := h.loadPlugin(dir)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

// InstallPlugin creates a symlink from the specified directory to the plugins
// directory
func (h *DefaultHandler) InstallPlugin(path string) error {
	// Validate the path to prevent path traversal
	absPath, err := validatePath(path)
	if err != nil {
		return fmt.Errorf("invalid plugin path: %v", err)
	}

	if _, err := os.Stat(h.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(h.Dir, 0755); err != nil {
			return err
		}
	}

	if err := h.validatePlugin(absPath); err != nil {
		return err
	}

	plugin, err := h.loadPlugin(absPath)
	if err != nil {
		return err
	}

	// Validate plugin name to prevent path traversal
	if !isValidPluginName(plugin.Name) {
		return fmt.Errorf("invalid plugin name: %s", plugin.Name)
	}

	targetPath := filepath.Join(h.Dir, plugin.Name)
	
	// Check if target already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("plugin with name '%s' already exists", plugin.Name)
	}

	return os.Symlink(plugin.Path, targetPath)
}

// UninstallPlugin removes the symlink from the plugin directory by plugin name
func (h *DefaultHandler) UninstallPlugin(name string) error {
	// Validate plugin name to prevent path traversal
	if !isValidPluginName(name) {
		return fmt.Errorf("invalid plugin name: %s", name)
	}

	plugins, err := h.GetPlugins()
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if plugin.Name == name {
			return os.Remove(filepath.Join(h.Dir, name))
		}
	}

	return fmt.Errorf("plugin '%s' was not found", name)
}

func (h *DefaultHandler) validatePlugin(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		return errors.New("plugin does not exist")
	}

	pluginFile := filepath.Join(dir, h.Filename)
	if _, err := os.Stat(pluginFile); err != nil {
		return fmt.Errorf("%s does not exist", h.Filename)
	}

	return nil
}

func (h *DefaultHandler) loadPlugin(dir string) (*Plugin, error) {
	// Validate the path to prevent path traversal
	absPath, err := validatePath(dir)
	if err != nil {
		return nil, fmt.Errorf("invalid plugin directory: %v", err)
	}

	// Use the new secure YAML loading function
	pluginYamlPath := filepath.Join(absPath, h.Filename)
	plugin, err := LoadPluginFromYAML(pluginYamlPath)
	if err != nil {
		return nil, err
	}

	// Set the path
	plugin.Path = absPath

	return plugin, nil
}

// LoadGoPlugin loads a cobra.Command from the Go plugin at the given path. The Go plugin must implement a command factory
// as follows:
//
// 	func New(settings *environment.Settings) *cobra.Command {
//		return &cobra.Command{...}
//  }
//
// If the file at the given path is not a Go plugin then the error, ErrNotAGoPlugin, is returned.
// If the Go plugin does not implement the command factory then an error is returned.
func (h *DefaultHandler) LoadGoPlugin(path string, settings *environment.Settings) (*cobra.Command, error) {
	// Validate the path to prevent path traversal
	absPath, err := validatePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid plugin path: %v", err)
	}

	p, err := goplugin.Open(absPath)
	if err != nil {
		return nil, ErrNotAGoPlugin
	}

	cmdFactorySymbol, err := p.Lookup(pluginFactoryMethod)
	if err != nil {
		return nil, fmt.Errorf("could not find symbol %s. Plugin must export this method", pluginFactoryMethod)
	}

	newCmd, ok := cmdFactorySymbol.(func(settings *environment.Settings) *cobra.Command)
	if !ok {
		return nil, fmt.Errorf("function %s does not match expected definition func(*environment.Settings) *cobra.Command", pluginFactoryMethod)
	}

	return newCmd(settings), nil
}

// validatePath ensures the path is valid and not potentially malicious
func validatePath(path string) (string, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", absPath)
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return "", ErrInvalidPluginPath
	}

	return absPath, nil
}

// isValidPluginName checks if the plugin name is valid and not potentially malicious
func isValidPluginName(name string) bool {
	// Check for path traversal attempts or special characters
	return !strings.Contains(name, "..") && 
		!strings.Contains(name, "/") && 
		!strings.Contains(name, "\\") &&
		!strings.Contains(name, " ") &&
		len(name) > 0
}

// validatePluginFields validates the plugin fields to ensure they are safe
func validatePluginFields(p *Plugin) error {
	if p.Name == "" {
		return errors.New("plugin name is required")
	}

	if !isValidPluginName(p.Name) {
		return fmt.Errorf("invalid plugin name: %s", p.Name)
	}

	if p.Command == nil {
		return errors.New("plugin command is required")
	}

	if p.Command.Base == "" {
		return errors.New("plugin command base is required")
	}

	return nil
}
