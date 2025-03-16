/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// DefaultConfigFilename is the default config filename
const DefaultConfigFilename = "config.yaml"

const contextTemplateString = `Network:	{{.Network}}
Organization:	{{.Organization}}
User:	{{.User}}
Channel:	{{.Channel}}
Orderers:
{{- range .Orderers}}
	{{.}}
{{- end}}
Peers:
{{- range .Peers}}
	{{.}}
{{- end}}`

const networkTemplateString = `Path:	{{.ConfigPath}}`

// Context contains network interaction parameters
type Context struct {
	Network      string   `yaml:",omitempty"`
	Organization string   `yaml:",omitempty"`
	User         string   `yaml:",omitempty"`
	Channel      string   `yaml:",omitempty"`
	Orderers     []string `yaml:",omitempty"`
	Peers        []string `yaml:",omitempty"`
}

func (c *Context) String() string {
	t := template.New("Context")
	data := new(bytes.Buffer)
	w := tabwriter.NewWriter(data, 4, 4, 4, ' ', 0)

	t.Parse(contextTemplateString)
	t.Execute(w, c)

	w.Flush()

	return data.String()
}

// Network contains a fabric network's configurations
type Network struct {
	// path to fabric go sdk config file
	ConfigPath string `yaml:"path,omitempty"`
}

func (n *Network) String() string {
	t := template.New("Network")
	data := new(bytes.Buffer)
	w := tabwriter.NewWriter(data, 4, 4, 4, ' ', 0)

	t.Parse(networkTemplateString)
	t.Execute(w, n)

	w.Flush()

	return data.String()
}

// Config contains information needed to manage fabric networks
type Config struct {
	Networks       map[string]*Network `yaml:",omitempty"`
	Contexts       map[string]*Context `yaml:",omitempty"`
	CurrentContext string              `yaml:"current-context,omitempty"`
}

// AddFlags appeneds config flags onto an existing flag set
func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.CurrentContext, "context", c.CurrentContext, "Override the current context")
}

// LoadFromFile populates config based on the specified path
func (c *Config) LoadFromFile(path string) error {
	// Validate the path to prevent path traversal
	if err := validateConfigPath(path); err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Use safe YAML unmarshaling
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	// Validate config after loading
	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

// Save writes the current config value to the specified path
func (c *Config) Save(path string) error {
	// Validate the path to prevent path traversal
	if err := validateConfigPath(path); err != nil {
		return err
	}

	// Validate config before saving
	if err := c.validate(); err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to a temporary file first, then rename to ensure atomic write
	tempFile := path + ".tmp"
	if err := ioutil.WriteFile(tempFile, data, 0600); err != nil {
		return err
	}

	if err := os.Rename(tempFile, path); err != nil {
		// Clean up the temporary file if rename fails
		os.Remove(tempFile)
		return err
	}

	return nil
}

// GetCurrentContext returns the current context
func (c *Config) GetCurrentContext() (*Context, error) {
	if c == nil || len(c.CurrentContext) == 0 {
		return nil, errors.New("current context is not set")
	}

	context, ok := c.Contexts[c.CurrentContext]
	if !ok {
		return nil, errors.New("current context does not exist")
	}

	return context, nil
}

// GetCurrentContextNetwork returns the current network
func (c *Config) GetCurrentContextNetwork() (*Network, error) {
	context, err := c.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	if len(context.Network) == 0 {
		return nil, errors.New("no network is set for current context")
	}

	network, ok := c.Networks[context.Network]
	if !ok {
		return nil, errors.New("network of the current context does not exist")
	}

	return network, nil
}

// validate performs validation on the config
func (c *Config) validate() error {
	// Validate networks
	for name, network := range c.Networks {
		if network == nil {
			return fmt.Errorf("network '%s' is nil", name)
		}
		
		// Validate network name
		if !isValidName(name) {
			return fmt.Errorf("invalid network name: %s", name)
		}
		
		// Validate config path
		if network.ConfigPath != "" && !isValidFilePath(network.ConfigPath) {
			return fmt.Errorf("invalid config path for network '%s': %s", name, network.ConfigPath)
		}
	}

	// Validate contexts
	for name, context := range c.Contexts {
		if context == nil {
			return fmt.Errorf("context '%s' is nil", name)
		}
		
		// Validate context name
		if !isValidName(name) {
			return fmt.Errorf("invalid context name: %s", name)
		}
		
		// Validate network reference
		if context.Network != "" && !isValidName(context.Network) {
			return fmt.Errorf("invalid network name in context '%s': %s", name, context.Network)
		}
	}

	// Validate current context
	if c.CurrentContext != "" && !isValidName(c.CurrentContext) {
		return fmt.Errorf("invalid current context name: %s", c.CurrentContext)
	}

	return nil
}

// NewConfig returns a new config
func NewConfig() *Config {
	return &Config{
		Networks: make(map[string]*Network),
		Contexts: make(map[string]*Context),
	}
}

// validateConfigPath validates the config file path
func validateConfigPath(path string) error {
	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return errors.New("path contains potentially unsafe '..' sequence")
	}

	// Ensure the path is absolute
	if !filepath.IsAbs(path) {
		return errors.New("config path must be absolute")
	}

	return nil
}

// isValidName checks if a name is valid
func isValidName(name string) bool {
	// Check for path traversal attempts or special characters
	return !strings.Contains(name, "..") && 
		!strings.Contains(name, "/") && 
		!strings.Contains(name, "\\") &&
		len(name) > 0
}

// isValidFilePath checks if a file path is valid
func isValidFilePath(path string) bool {
	// Check for path traversal attempts
	return !strings.Contains(path, "..")
}
