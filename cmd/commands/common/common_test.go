package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleCollectionsConfigGood = `[
	{
		"name": "foo",
		"policy": "OR('A.member', 'B.member')",
		"requiredPeerCount": 3,
		"maxPeerCount": 483279847,
		"blockToLive":10,
		"memberOnlyRead": true,
		"memberOnlyWrite": true
	}
]`

const sampleCollectionsConfigBad = `[
	{
		"name": "foo",
		"policy": "barf",
		"requiredPeerCount": 3,
		"maxPeerCount": 483279847
	}
]`

func TestGetCollectionsConfigFromBytes(t *testing.T) {
	config, err := GetCollectionsConfigFromBytes([]byte(sampleCollectionsConfigGood))
	assert.Nil(t, err)
	assert.NotNil(t, config)
}

func TestGetCollectionsConfigFromBytesError(t *testing.T) {
	config, err := GetCollectionsConfigFromBytes([]byte(sampleCollectionsConfigBad))
	assert.NotNil(t, err)
	assert.Nil(t, config)
}

func TestGetChaincodePolicy(t *testing.T) {
	policy, err := GetChaincodePolicy("OR('MSP.member', 'MSP.WITH.DOTS.member', 'MSP-WITH-DASHES.member')")
	assert.Nil(t, err)
	assert.NotNil(t, policy)
}

func TestGetChaincodePolicyError(t *testing.T) {
	policy, err := GetChaincodePolicy("NOT A VALID POLICY)")
	assert.NotNil(t, err)
	assert.Nil(t, policy)
}
