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

	Context         *Context
	CryptoConfig    string
	CredentialStore string

	Channels      map[string]*Channel
	Organizations map[string]*Organization
	Orderers      map[string]*Orderer
	Peers         map[string]*Peer
}

// Context contains the active focus of the profile
type Context struct {
	Organization string
	Identity     string
	Channel      string
	Orderers     []string
	Peers        []string
}

// Channel contains configuration details for a channel
type Channel struct {
	ID    string
	Peers []string
}

// Organization contains configuration details for an organization
type Organization struct {
	ID    string
	MSP   *MSP
	Peers []string
}

// MSP contains configuration details for a msp
type MSP struct {
	ID    string
	Store string
}

// Orderer contains configuration details for a orderer
type Orderer struct {
	ID  string
	URL string
	TLS string
}

// Peer contains configuration details for a peer
type Peer struct {
	ID       string
	URL      string
	EventURL string
	TLS      string

	ChannelOptions map[string]interface{}
	GRPCOptions    map[string]interface{}
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
