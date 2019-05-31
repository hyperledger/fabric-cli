/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"os"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type factory struct {
	config  string
	context *environment.Context
}

// interface implementation check
var _ Factory = &factory{}

// NewFactory creates a factory for the given profile/context
func NewFactory(config *environment.Config) (Factory, error) {
	context, err := config.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	network, err := config.GetCurrentContextNetwork()
	if err != nil {
		return nil, err
	}

	return &factory{
		config:  network.ConfigPath,
		context: context,
	}, nil
}

func (f *factory) SDK() (SDK, error) {
	sdk, err := fabsdk.New(config.FromFile(os.ExpandEnv(f.config)))
	if err != nil {
		return nil, err
	}

	return sdk, nil
}

func (f *factory) Channel() (Channel, error) {
	sdk, err := f.SDK()
	if err != nil {
		return nil, err
	}

	ctx := sdk.ChannelContext(
		f.context.Channel,
		fabsdk.WithUser(f.context.User),
		fabsdk.WithOrg(f.context.Organization),
	)

	client, err := channel.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *factory) Event() (Event, error) {
	sdk, err := f.SDK()
	if err != nil {
		return nil, err
	}

	ctx := sdk.ChannelContext(
		f.context.Channel,
		fabsdk.WithUser(f.context.User),
		fabsdk.WithOrg(f.context.Organization),
	)

	client, err := event.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *factory) Ledger() (Ledger, error) {
	sdk, err := f.SDK()
	if err != nil {
		return nil, err
	}

	ctx := sdk.ChannelContext(
		f.context.Channel,
		fabsdk.WithUser(f.context.User),
		fabsdk.WithOrg(f.context.Organization),
	)

	client, err := ledger.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *factory) ResourceManagement() (ResourceManagement, error) {
	sdk, err := f.SDK()
	if err != nil {
		return nil, err
	}

	ctx := sdk.Context(fabsdk.WithUser(f.context.User),
		fabsdk.WithOrg(f.context.Organization))

	client, err := resmgmt.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *factory) MSP() (MSP, error) {
	sdk, err := f.SDK()
	if err != nil {
		return nil, err
	}

	ctx := sdk.Context(fabsdk.WithUser(f.context.User),
		fabsdk.WithOrg(f.context.Organization))

	client, err := msp.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
