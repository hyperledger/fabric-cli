package chaincode

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

type collectionConfigJSON struct {
	Name            string `json:"name"`
	Policy          string `json:"policy"`
	RequiredCount   int32  `json:"requiredPeerCount"`
	MaxPeerCount    int32  `json:"maxPeerCount"`
	BlockToLive     uint64 `json:"blockToLive"`
	MemberOnlyRead  bool   `json:"memberOnlyRead"`
	MemberOnlyWrite bool   `json:"memberOnlyWrite"`
}

func getCollectionConfigFromFile(path string) ([]*common.CollectionConfig, error) {
	if len(path) == 0 {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("error reading collections config file")
	}

	return getCollectionsConfigFromBytes(bytes)
}

func getCollectionsConfigFromBytes(bytes []byte) ([]*common.CollectionConfig, error) {
	var cconf []collectionConfigJSON
	if err := json.Unmarshal(bytes, &cconf); err != nil {
		return nil, errors.New("error unmarshalling collections config")
	}

	ccarray := make([]*common.CollectionConfig, 0, len(cconf))
	for _, cconfitem := range cconf {
		p, err := cauthdsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, err
		}
		cpc := &common.CollectionPolicyConfig{
			Payload: &common.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}
		cc := &common.CollectionConfig{
			Payload: &common.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &common.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: cconfitem.RequiredCount,
					MaximumPeerCount:  cconfitem.MaxPeerCount,
				},
			},
		}
		ccarray = append(ccarray, cc)
	}
	return ccarray, nil
}

func getChaincodePolicy(policyString string) (*common.SignaturePolicyEnvelope, error) {
	if len(policyString) == 0 {
		return cauthdsl.AcceptAllPolicy, nil
	}

	policy, err := cauthdsl.FromString(policyString)
	if err != nil {
		return nil, errors.New("error parsing chaincode policy")
	}
	return policy, nil
}
