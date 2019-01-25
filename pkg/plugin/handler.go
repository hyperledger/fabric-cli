/*
Copyright © 2019 State Street Bank and Trust Company.  All rights reserved

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Handler defines the required actions for managing plugins
type Handler interface {
	GetPlugins() ([]*Plugin, error)
	InstallPlugin(path string) error
	UninstallPlugin(name string) error
}

// DefaultHandler is the default plugin handler
type DefaultHandler struct {
	Dir              string
	MetadataFileName string
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
	if _, err := os.Stat(h.Dir); err != nil {
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

	if _, err := os.Stat(filepath.Join(dir, h.MetadataFileName)); err != nil {
		return fmt.Errorf("%s does not exist", h.MetadataFileName)
	}

	return nil
}

func (h *DefaultHandler) loadPlugin(dir string) (*Plugin, error) {
	p, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(p, h.MetadataFileName))
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
