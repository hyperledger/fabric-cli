/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"bytes"
	"text/template"
)

// Profile contains metadata for a fabric network
type Profile struct {
	Name string

	Context         *Context `yaml:",omitempty"`
	CryptoConfig    string   `yaml:",omitempty"`
	CredentialStore string   `yaml:",omitempty"`

	Channels      map[string]*Channel      `yaml:",omitempty"`
	Organizations map[string]*Organization `yaml:",omitempty"`
	Orderers      map[string]*Orderer      `yaml:",omitempty"`
	Peers         map[string]*Peer         `yaml:",omitempty"`
}

// Context contains the active focus of the profile
type Context struct {
	Organization string   `yaml:",omitempty"`
	Identity     string   `yaml:",omitempty"`
	Channel      string   `yaml:",omitempty"`
	Orderers     []string `yaml:",omitempty"`
	Peers        []string `yaml:",omitempty"`
}

// Channel contains configuration details for a channel
type Channel struct {
	ID    string   `yaml:",omitempty"`
	Peers []string `yaml:",omitempty"`
}

// Organization contains configuration details for an organization
type Organization struct {
	ID    string   `yaml:",omitempty"`
	MSP   *MSP     `yaml:",omitempty"`
	Peers []string `yaml:",omitempty"`
}

// MSP contains configuration details for a msp
type MSP struct {
	ID    string `yaml:",omitempty"`
	Store string `yaml:",omitempty"`
}

// Orderer contains configuration details for a orderer
type Orderer struct {
	ID  string `yaml:",omitempty"`
	URL string `yaml:",omitempty"`
	TLS string `yaml:",omitempty"`
}

// Peer contains configuration details for a peer
type Peer struct {
	ID       string `yaml:",omitempty"`
	URL      string `yaml:",omitempty"`
	EventURL string `yaml:",omitempty"`
	TLS      string `yaml:",omitempty"`

	ChannelOptions map[string]interface{} `yaml:",omitempty"`
	GRPCOptions    map[string]interface{} `yaml:",omitempty"`
}

// ToTemplate transforms the profile into the provided template
func (p *Profile) ToTemplate(path string) ([]byte, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	if err := tmpl.Execute(buffer, p); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
