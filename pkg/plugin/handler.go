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
	goplugin "plugin"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const pluginFactoryMethod = "New"

// ErrNotAGoPlugin is returned when attempting to load a file that's not a Go plugin
var ErrNotAGoPlugin = errors.New("not a Go plugin")

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
	if _, err := os.Stat(h.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(h.Dir, 0755); err != nil {
			return err
		}
	}

	if err := h.validatePlugin(path); err != nil {
		return err
	}

	plugin, err := h.loadPlugin(path)
	if err != nil {
		return err
	}

	return os.Symlink(plugin.Path, filepath.Join(h.Dir, plugin.Name))
}

// UninstallPlugin removes the symlink from the plugin directory by plugin name
func (h *DefaultHandler) UninstallPlugin(name string) error {
	plugins, err := h.GetPlugins()
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if plugin.Name == name {
			return os.Remove(plugin.Path)
		}
	}

	return fmt.Errorf("plugin '%s' was not found", name)
}

func (h *DefaultHandler) validatePlugin(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		return errors.New("plugin does not exist")
	}

	if _, err := os.Stat(filepath.Join(dir, h.Filename)); err != nil {
		return fmt.Errorf("%s does not exist", h.Filename)
	}

	return nil
}

func (h *DefaultHandler) loadPlugin(dir string) (*Plugin, error) {
	p, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(p, h.Filename))
	if err != nil {
		return nil, err
	}

	plugin := &Plugin{
		Path: p,
	}

	if err := yaml.Unmarshal(data, &plugin); err != nil {
		return nil, err
	}

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
	p, err := goplugin.Open(path)
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
