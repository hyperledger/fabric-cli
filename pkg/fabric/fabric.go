/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"github.com/hyperledger/fabric-cli/pkg/environment"
)

// client is a base fabric client
type client struct {
	factory Factory
}

// newClient creates a base fabric client that populates the default factory and
// applies options.
func newSDKClient(profile *environment.Profile, options []Option) *client {
	c := &client{
		factory: NewFactory(profile),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// ChannelClient encapsulates the SDK channel client
type ChannelClient struct {
	*client
	Channel
}

// NewChannelClient returns a new channel client
func NewChannelClient(profile *environment.Profile, options ...Option) (*ChannelClient, error) {
	c := &ChannelClient{
		client: newSDKClient(profile, options),
	}

	client, err := c.factory.Channel()
	if err != nil {
		return nil, err
	}

	c.Channel = client

	return c, nil
}

// EventClient encapsulates the SDK event client
type EventClient struct {
	*client
	Event
}

// NewEventClient returns a new event client
func NewEventClient(profile *environment.Profile, options ...Option) (*EventClient, error) {
	c := &EventClient{
		client: newSDKClient(profile, options),
	}

	client, err := c.factory.Event()
	if err != nil {
		return nil, err
	}

	c.Event = client

	return c, nil
}

// LedgerClient encapsulates the SDK ledger client
type LedgerClient struct {
	*client
	Ledger
}

// NewLedgerClient returns a new ledger client
func NewLedgerClient(profile *environment.Profile, options ...Option) (*LedgerClient, error) {
	c := &LedgerClient{
		client: newSDKClient(profile, options),
	}

	client, err := c.factory.Ledger()
	if err != nil {
		return nil, err
	}

	c.Ledger = client

	return c, nil
}

// ResourceManagementClient encapsulates the SDK resmgmt client
type ResourceManagementClient struct {
	*client
	ResourceManagement
}

// NewResourceManagementClient returns a new resource management client
func NewResourceManagementClient(profile *environment.Profile, options ...Option) (*ResourceManagementClient, error) {
	c := &ResourceManagementClient{
		client: newSDKClient(profile, options),
	}

	client, err := c.factory.ResourceManagement()
	if err != nil {
		return nil, err
	}

	c.ResourceManagement = client

	return c, nil
}

// MSPClient encapsulates the SDK msp client
type MSPClient struct {
	*client
	MSP
}

// NewMSPClient returns a new msp client
func NewMSPClient(profile *environment.Profile, options ...Option) (*MSPClient, error) {
	c := &MSPClient{
		client: newSDKClient(profile, options),
	}

	client, err := c.factory.MSP()
	if err != nil {
		return nil, err
	}

	c.MSP = client

	return c, nil
}
