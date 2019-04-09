/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"bytes"
	"text/template"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric/templates"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type factory struct {
	profile *environment.Profile
}

// interface implementation check
var _ Factory = &factory{}

// NewFactory creates a factory for the given profile/context
func NewFactory(profile *environment.Profile) Factory {
	return &factory{
		profile: profile,
	}
}

func (f *factory) SDK() (SDK, error) {
	tmpl, err := template.New("config").Parse(templates.Config)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	if err := tmpl.Execute(buffer, f.profile); err != nil {
		return nil, err
	}

	sdk, err := fabsdk.New(config.FromRaw(buffer.Bytes(), "yaml"))
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

	ctx := sdk.ChannelContext(f.profile.Context.Channel)

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

	ctx := sdk.ChannelContext(f.profile.Context.Channel)

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

	ctx := sdk.ChannelContext(f.profile.Context.Channel)

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

	ctx := sdk.Context(fabsdk.WithUser(f.profile.Context.Identity),
		fabsdk.WithOrg(f.profile.Context.Organization))

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

	ctx := sdk.Context(fabsdk.WithUser(f.profile.Context.Identity),
		fabsdk.WithOrg(f.profile.Context.Organization))

	client, err := msp.New(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
