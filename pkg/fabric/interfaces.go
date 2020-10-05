/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Factory provides abstractions for the various SDK clients
type Factory interface {
	SDK() (SDK, error)
	Channel() (Channel, error)
	Event() (Event, error)
	Ledger() (Ledger, error)
	ResourceManagement() (ResourceManagement, error)
	MSP() (MSP, error)
}

// SDK defines the context methods for the various SDK clients
type SDK interface {
	ChannelContext(channelID string, options ...fabsdk.ContextOption) context.ChannelProvider
	Context(options ...fabsdk.ContextOption) context.ClientProvider
	Config() (core.ConfigBackend, error)
	CloseContext(ctxt fab.ClientContext)
	Close()
}

// Channel defines the methods implemented by SDK channel client
type Channel interface {
	Execute(request channel.Request, options ...channel.RequestOption) (channel.Response, error)
	InvokeHandler(handler invoke.Handler, request channel.Request, options ...channel.RequestOption) (channel.Response, error)
	Query(request channel.Request, options ...channel.RequestOption) (channel.Response, error)
	RegisterChaincodeEvent(chainCodeID string, eventFilter string) (fab.Registration, <-chan *fab.CCEvent, error)
	UnregisterChaincodeEvent(registration fab.Registration)
}

// Event defines the methods implemented by SDK event client
type Event interface {
	RegisterBlockEvent(filter ...fab.BlockFilter) (fab.Registration, <-chan *fab.BlockEvent, error)
	RegisterChaincodeEvent(ccID, eventFilter string) (fab.Registration, <-chan *fab.CCEvent, error)
	RegisterFilteredBlockEvent() (fab.Registration, <-chan *fab.FilteredBlockEvent, error)
	RegisterTxStatusEvent(txID string) (fab.Registration, <-chan *fab.TxStatusEvent, error)
	Unregister(reg fab.Registration)
}

// Ledger defines the methods implemented by SDK ledger client
type Ledger interface {
	QueryBlock(blockNumber uint64, options ...ledger.RequestOption) (*common.Block, error)
	QueryBlockByHash(blockHash []byte, options ...ledger.RequestOption) (*common.Block, error)
	QueryBlockByTxID(txID fab.TransactionID, options ...ledger.RequestOption) (*common.Block, error)
	QueryConfig(options ...ledger.RequestOption) (fab.ChannelCfg, error)
	QueryInfo(options ...ledger.RequestOption) (*fab.BlockchainInfoResponse, error)
	QueryTransaction(transactionID fab.TransactionID, options ...ledger.RequestOption) (*pb.ProcessedTransaction, error)
}

// ResourceManagement defines the methods implemented by SDK resmgmt client
type ResourceManagement interface {
	CreateConfigSignature(signer mspctx.SigningIdentity, channelConfigPath string) (*common.ConfigSignature, error)
	CreateConfigSignatureData(signer mspctx.SigningIdentity, channelConfigPath string) (signatureHeaderData resource.ConfigSignatureData, e error)
	InstallCC(req resmgmt.InstallCCRequest, options ...resmgmt.RequestOption) ([]resmgmt.InstallCCResponse, error)
	InstantiateCC(channelID string, req resmgmt.InstantiateCCRequest, options ...resmgmt.RequestOption) (resmgmt.InstantiateCCResponse, error)
	JoinChannel(channelID string, options ...resmgmt.RequestOption) error
	QueryChannels(options ...resmgmt.RequestOption) (*pb.ChannelQueryResponse, error)
	QueryCollectionsConfig(channelID string, chaincodeName string, options ...resmgmt.RequestOption) (*pb.CollectionConfigPackage, error)
	QueryConfigFromOrderer(channelID string, options ...resmgmt.RequestOption) (fab.ChannelCfg, error)
	QueryInstalledChaincodes(options ...resmgmt.RequestOption) (*pb.ChaincodeQueryResponse, error)
	QueryInstantiatedChaincodes(channelID string, options ...resmgmt.RequestOption) (*pb.ChaincodeQueryResponse, error)
	SaveChannel(req resmgmt.SaveChannelRequest, options ...resmgmt.RequestOption) (resmgmt.SaveChannelResponse, error)
	UpgradeCC(channelID string, req resmgmt.UpgradeCCRequest, options ...resmgmt.RequestOption) (resmgmt.UpgradeCCResponse, error)
	LifecycleInstallCC(req resmgmt.LifecycleInstallCCRequest, options ...resmgmt.RequestOption) ([]resmgmt.LifecycleInstallCCResponse, error)
	LifecycleApproveCC(channelID string, req resmgmt.LifecycleApproveCCRequest, options ...resmgmt.RequestOption) (fab.TransactionID, error)
	LifecycleCommitCC(channelID string, req resmgmt.LifecycleCommitCCRequest, options ...resmgmt.RequestOption) (fab.TransactionID, error)
	LifecycleQueryInstalledCC(options ...resmgmt.RequestOption) ([]resmgmt.LifecycleInstalledCC, error)
	LifecycleQueryApprovedCC(channelID string, req resmgmt.LifecycleQueryApprovedCCRequest,
		options ...resmgmt.RequestOption) (resmgmt.LifecycleApprovedChaincodeDefinition, error)
	LifecycleCheckCCCommitReadiness(channelID string, req resmgmt.LifecycleCheckCCCommitReadinessRequest,
		options ...resmgmt.RequestOption) (resmgmt.LifecycleCheckCCCommitReadinessResponse, error)
	LifecycleQueryCommittedCC(channelID string, req resmgmt.LifecycleQueryCommittedCCRequest,
		options ...resmgmt.RequestOption) ([]resmgmt.LifecycleChaincodeDefinition, error)
}

// MSP defines the methods implemented by SDK msp client
type MSP interface {
	AddAffiliation(request *msp.AffiliationRequest) (*msp.AffiliationResponse, error)
	CreateIdentity(request *msp.IdentityRequest) (*msp.IdentityResponse, error)
	CreateSigningIdentity(opts ...mspctx.SigningIdentityOption) (mspctx.SigningIdentity, error)
	Enroll(enrollmentID string, opts ...msp.EnrollmentOption) error
	GetAffiliation(affiliation string, options ...msp.RequestOption) (*msp.AffiliationResponse, error)
	GetAllAffiliations(options ...msp.RequestOption) (*msp.AffiliationResponse, error)
	GetAllIdentities(options ...msp.RequestOption) ([]*msp.IdentityResponse, error)
	GetCAInfo() (*msp.GetCAInfoResponse, error)
	GetIdentity(ID string, options ...msp.RequestOption) (*msp.IdentityResponse, error)
	GetSigningIdentity(id string) (mspctx.SigningIdentity, error)
	ModifyAffiliation(request *msp.ModifyAffiliationRequest) (*msp.AffiliationResponse, error)
	ModifyIdentity(request *msp.IdentityRequest) (*msp.IdentityResponse, error)
	Reenroll(enrollmentID string, opts ...msp.EnrollmentOption) error
	Register(request *msp.RegistrationRequest) (string, error)
	RemoveAffiliation(request *msp.AffiliationRequest) (*msp.AffiliationResponse, error)
	RemoveIdentity(request *msp.RemoveIdentityRequest) (*msp.IdentityResponse, error)
	Revoke(request *msp.RevocationRequest) (*msp.RevocationResponse, error)
}
