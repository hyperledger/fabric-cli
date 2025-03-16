/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hyperledger/fabric-cli/cmd/commands"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/plugin"
)

// NewFabricCommand returns a new root command for fabric
func NewFabricCommand(settings *environment.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fabric",
		Short: "The command line interface for Hyperledger Fabric",
	}

	// load all built in commands into the root command
	cmd.AddCommand(commands.All(settings)...)

	cmd.SetOutput(settings.Streams.Out)

	return cmd
}

// NewDefaultFabricCommand returns a new default root commad for fabric
func NewDefaultFabricCommand(settings *environment.Settings, args []string) *cobra.Command {
	cmd := NewFabricCommand(settings)
	flags := cmd.PersistentFlags()

	settings.AddFlags(flags)
	flags.Parse(args)

	if err := settings.Init(flags); err != nil {
		fmt.Fprintf(settings.Streams.Err, "An error occurred while loading configurations: %v\n", err)
		os.Exit(1)
	}

	// must resolve home and config file before loading config flags
	settings.Config.AddFlags(flags)

	if err := loadPlugins(cmd, settings, &plugin.DefaultHandler{
		Dir:      settings.Home.Plugins(),
		Filename: plugin.DefaultFilename,
	}); err != nil {
		fmt.Fprintf(settings.Streams.Err, "An error occurred while loading plugins: %v\n", err)
		os.Exit(1)
	}

	return cmd
}

func main() {
	settings := environment.NewDefaultSettings()
	cmd := NewDefaultFabricCommand(settings, os.Args[1:])

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

// loadPlugins processes all of the installed plugins, wraps them with cobra,
// and adds them to the root command
func loadPlugins(cmd *cobra.Command, settings *environment.Settings, handler plugin.Handler) error {
	if settings.DisablePlugins {
		return nil
	}

	plugins, err := handler.GetPlugins()
	if err != nil {
		return err
	}

	settings.SetupPluginEnvironment()

	for _, p := range plugins {
		c, err := loadPlugin(p, settings, handler)
		if err != nil {
			return err
		}
		cmd.AddCommand(c)
	}

	return nil
}

// loadPlugin loads the given plugin as either a Go plugin or a wrapped executable
func loadPlugin(p *plugin.Plugin, settings *environment.Settings, handler plugin.Handler) (*cobra.Command, error) {
	// Validate plugin path to prevent path traversal
	path, err := validatePluginPath(os.ExpandEnv(p.Command.Base), settings.Home.Plugins())
	if err != nil {
		return nil, fmt.Errorf("invalid plugin path: %v", err)
	}

	c, err := handler.LoadGoPlugin(path, settings)
	if err == nil {
		return c, nil
	}
	if err != plugin.ErrNotAGoPlugin {
		return nil, err
	}

	// Validate allowed commands
	if !isAllowedCommand(path) {
		return nil, fmt.Errorf("command not in allowed list: %s", path)
	}

	return &cobra.Command{
		Use:   p.Name,
		Short: p.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Sanitize arguments to prevent command injection
			sanitizedArgs := sanitizeArgs(append(p.Command.Args, args...))
			e := exec.Command(path, sanitizedArgs...)
			e.Env = os.Environ()
			e.Stdin = settings.Streams.In
			e.Stdout = settings.Streams.Out
			e.Stderr = settings.Streams.Err
			return e.Run()
		},
	}, nil
}

// validatePluginPath ensures the plugin path is within the allowed plugins directory
func validatePluginPath(path, pluginsDir string) (string, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("plugin path does not exist: %s", absPath)
	}

	// For security, only allow plugins from specific directories
	absPluginsDir, err := filepath.Abs(pluginsDir)
	if err != nil {
		return "", err
	}

	// Check if the plugin is within the plugins directory or is in a system path
	if !strings.HasPrefix(absPath, absPluginsDir) && !isInSystemPath(absPath) {
		return "", fmt.Errorf("plugin must be in the plugins directory or a system path")
	}

	return absPath, nil
}

// isInSystemPath checks if the given path is in a system path
func isInSystemPath(path string) bool {
	// Get system PATH
	systemPaths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	
	// Check if the path is in any of the system paths
	for _, sysPath := range systemPaths {
		absPath, err := filepath.Abs(sysPath)
		if err != nil {
			continue
		}
		if strings.HasPrefix(path, absPath) {
			return true
		}
	}
	
	return false
}

// isAllowedCommand checks if the command is in the allowed list
func isAllowedCommand(cmd string) bool {
	// Define a whitelist of allowed commands
	// This should be configured based on your security requirements
	allowedCommands := []string{
		"cryptogen",
		"configtxgen",
		"configtxlator",
		"peer",
		"orderer",
		"discover",
		"idemixgen",
	}

	cmdBase := filepath.Base(cmd)
	for _, allowed := range allowedCommands {
		if cmdBase == allowed {
			return true
		}
	}

	return false
}

// sanitizeArgs sanitizes command arguments to prevent command injection
func sanitizeArgs(args []string) []string {
	sanitized := make([]string, len(args))
	for i, arg := range args {
		// Remove any characters that could be used for command injection
		// This is a basic implementation - you may need more sophisticated sanitization
		sanitized[i] = strings.ReplaceAll(arg, ";", "")
		sanitized[i] = strings.ReplaceAll(sanitized[i], "|", "")
		sanitized[i] = strings.ReplaceAll(sanitized[i], "&", "")
		sanitized[i] = strings.ReplaceAll(sanitized[i], ">", "")
		sanitized[i] = strings.ReplaceAll(sanitized[i], "<", "")
	}
	return sanitized
}
