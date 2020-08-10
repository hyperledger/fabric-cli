package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)

// CollectionConfigJSON contains the parameters for a collection configuration
type CollectionConfigJSON struct {
	Name            string `json:"name"`
	Policy          string `json:"policy"`
	RequiredCount   int32  `json:"requiredPeerCount"`
	MaxPeerCount    int32  `json:"maxPeerCount"`
	BlockToLive     uint64 `json:"blockToLive"`
	MemberOnlyRead  bool   `json:"memberOnlyRead"`
	MemberOnlyWrite bool   `json:"memberOnlyWrite"`
}

// GetCollectionConfigFromFile returns the collection config from the given file
func GetCollectionConfigFromFile(path string) ([]*pb.CollectionConfig, error) {
	if len(path) == 0 {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("error reading collections config file")
	}

	return GetCollectionsConfigFromBytes(bytes)
}

// GetCollectionsConfigFromBytes returns the collection config from the given byte array
func GetCollectionsConfigFromBytes(bytes []byte) ([]*pb.CollectionConfig, error) {
	var cconf []CollectionConfigJSON
	if err := json.Unmarshal(bytes, &cconf); err != nil {
		return nil, errors.New("error unmarshalling collections config")
	}

	ccarray := make([]*pb.CollectionConfig, 0, len(cconf))
	for _, cconfitem := range cconf {
		p, err := policydsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, err
		}
		cpc := &pb.CollectionPolicyConfig{
			Payload: &pb.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}
		cc := &pb.CollectionConfig{
			Payload: &pb.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &pb.StaticCollectionConfig{
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

// GetChaincodePolicy returns the signature policy from the given policy string
func GetChaincodePolicy(policyString string) (*common.SignaturePolicyEnvelope, error) {
	if len(policyString) == 0 {
		return policydsl.AcceptAllPolicy, nil
	}

	policy, err := policydsl.FromString(policyString)
	if err != nil {
		return nil, errors.New("error parsing chaincode policy")
	}
	return policy, nil
}
