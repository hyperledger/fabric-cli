/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
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
	fs.StringVar(&c.CurrentContext, "context", c.CurrentContext, "override the current context")
}

// LoadFromFile populates config based on the specified path
func (c *Config) LoadFromFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	return nil
}

// Save writes the current config value to the specified path
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
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

// NewConfig returns a new config
func NewConfig() *Config {
	return &Config{
		Networks: make(map[string]*Network),
		Contexts: make(map[string]*Context),
	}
}
