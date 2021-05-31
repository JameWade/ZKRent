// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/state/types.proto

package state

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	types "github.com/ChengtayChain/ChengtayChain/abci/types"
	types1 "github.com/ChengtayChain/ChengtayChain/proto/types"
	version "github.com/ChengtayChain/ChengtayChain/proto/version"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ABCIResponses retains the responses
// of the various ABCI calls during block processing.
// It is persisted to disk for each height before calling Commit.
type ABCIResponses struct {
	DeliverTxs           []*types.ResponseDeliverTx `protobuf:"bytes,1,rep,name=deliver_txs,json=deliverTxs,proto3" json:"deliver_txs,omitempty"`
	EndBlock             *types.ResponseEndBlock    `protobuf:"bytes,2,opt,name=end_block,json=endBlock,proto3" json:"end_block,omitempty"`
	BeginBlock           *types.ResponseBeginBlock  `protobuf:"bytes,3,opt,name=begin_block,json=beginBlock,proto3" json:"begin_block,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ABCIResponses) Reset()         { *m = ABCIResponses{} }
func (m *ABCIResponses) String() string { return proto.CompactTextString(m) }
func (*ABCIResponses) ProtoMessage()    {}
func (*ABCIResponses) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{0}
}
func (m *ABCIResponses) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ABCIResponses.Unmarshal(m, b)
}
func (m *ABCIResponses) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ABCIResponses.Marshal(b, m, deterministic)
}
func (m *ABCIResponses) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ABCIResponses.Merge(m, src)
}
func (m *ABCIResponses) XXX_Size() int {
	return xxx_messageInfo_ABCIResponses.Size(m)
}
func (m *ABCIResponses) XXX_DiscardUnknown() {
	xxx_messageInfo_ABCIResponses.DiscardUnknown(m)
}

var xxx_messageInfo_ABCIResponses proto.InternalMessageInfo

func (m *ABCIResponses) GetDeliverTxs() []*types.ResponseDeliverTx {
	if m != nil {
		return m.DeliverTxs
	}
	return nil
}

func (m *ABCIResponses) GetEndBlock() *types.ResponseEndBlock {
	if m != nil {
		return m.EndBlock
	}
	return nil
}

func (m *ABCIResponses) GetBeginBlock() *types.ResponseBeginBlock {
	if m != nil {
		return m.BeginBlock
	}
	return nil
}

// ValidatorsInfo represents the latest validator set, or the last height it changed
type ValidatorsInfo struct {
	ValidatorSet         *types1.ValidatorSet `protobuf:"bytes,1,opt,name=validator_set,json=validatorSet,proto3" json:"validator_set,omitempty"`
	LastHeightChanged    int64                `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ValidatorsInfo) Reset()         { *m = ValidatorsInfo{} }
func (m *ValidatorsInfo) String() string { return proto.CompactTextString(m) }
func (*ValidatorsInfo) ProtoMessage()    {}
func (*ValidatorsInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{1}
}
func (m *ValidatorsInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidatorsInfo.Unmarshal(m, b)
}
func (m *ValidatorsInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidatorsInfo.Marshal(b, m, deterministic)
}
func (m *ValidatorsInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorsInfo.Merge(m, src)
}
func (m *ValidatorsInfo) XXX_Size() int {
	return xxx_messageInfo_ValidatorsInfo.Size(m)
}
func (m *ValidatorsInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorsInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorsInfo proto.InternalMessageInfo

func (m *ValidatorsInfo) GetValidatorSet() *types1.ValidatorSet {
	if m != nil {
		return m.ValidatorSet
	}
	return nil
}

func (m *ValidatorsInfo) GetLastHeightChanged() int64 {
	if m != nil {
		return m.LastHeightChanged
	}
	return 0
}

// ConsensusParamsInfo represents the latest consensus params, or the last height it changed
type ConsensusParamsInfo struct {
	ConsensusParams      types1.ConsensusParams `protobuf:"bytes,1,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params"`
	LastHeightChanged    int64                  `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ConsensusParamsInfo) Reset()         { *m = ConsensusParamsInfo{} }
func (m *ConsensusParamsInfo) String() string { return proto.CompactTextString(m) }
func (*ConsensusParamsInfo) ProtoMessage()    {}
func (*ConsensusParamsInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{2}
}
func (m *ConsensusParamsInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusParamsInfo.Unmarshal(m, b)
}
func (m *ConsensusParamsInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusParamsInfo.Marshal(b, m, deterministic)
}
func (m *ConsensusParamsInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusParamsInfo.Merge(m, src)
}
func (m *ConsensusParamsInfo) XXX_Size() int {
	return xxx_messageInfo_ConsensusParamsInfo.Size(m)
}
func (m *ConsensusParamsInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusParamsInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusParamsInfo proto.InternalMessageInfo

func (m *ConsensusParamsInfo) GetConsensusParams() types1.ConsensusParams {
	if m != nil {
		return m.ConsensusParams
	}
	return types1.ConsensusParams{}
}

func (m *ConsensusParamsInfo) GetLastHeightChanged() int64 {
	if m != nil {
		return m.LastHeightChanged
	}
	return 0
}

type Version struct {
	Consensus            version.Consensus `protobuf:"bytes,1,opt,name=consensus,proto3" json:"consensus"`
	Software             string            `protobuf:"bytes,2,opt,name=software,proto3" json:"software,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{3}
}
func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetConsensus() version.Consensus {
	if m != nil {
		return m.Consensus
	}
	return version.Consensus{}
}

func (m *Version) GetSoftware() string {
	if m != nil {
		return m.Software
	}
	return ""
}

type State struct {
	Version Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version"`
	// immutable
	ChainID string `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// LastBlockHeight=0 at genesis (ie. block(H=0) does not exist)
	LastBlockHeight int64          `protobuf:"varint,3,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	LastBlockID     types1.BlockID `protobuf:"bytes,4,opt,name=last_block_id,json=lastBlockId,proto3" json:"last_block_id"`
	LastBlockTime   time.Time      `protobuf:"bytes,5,opt,name=last_block_time,json=lastBlockTime,proto3,stdtime" json:"last_block_time"`
	// LastValidators is used to validate block.LastCommit.
	// Validators are persisted to the database separately every time they change,
	// so we can query for historical validator sets.
	// Note that if s.LastBlockHeight causes a valset change,
	// we set s.LastHeightValidatorsChanged = s.LastBlockHeight + 1 + 1
	// Extra +1 due to nextValSet delay.
	NextValidators              *types1.ValidatorSet `protobuf:"bytes,6,opt,name=next_validators,json=nextValidators,proto3" json:"next_validators,omitempty"`
	Validators                  *types1.ValidatorSet `protobuf:"bytes,7,opt,name=validators,proto3" json:"validators,omitempty"`
	LastValidators              *types1.ValidatorSet `protobuf:"bytes,8,opt,name=last_validators,json=lastValidators,proto3" json:"last_validators,omitempty"`
	LastHeightValidatorsChanged int64                `protobuf:"varint,9,opt,name=last_height_validators_changed,json=lastHeightValidatorsChanged,proto3" json:"last_height_validators_changed,omitempty"`
	// Consensus parameters used for validating blocks.
	// Changes returned by EndBlock and updated after Commit.
	ConsensusParams                  types1.ConsensusParams `protobuf:"bytes,10,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params"`
	LastHeightConsensusParamsChanged int64                  `protobuf:"varint,11,opt,name=last_height_consensus_params_changed,json=lastHeightConsensusParamsChanged,proto3" json:"last_height_consensus_params_changed,omitempty"`
	// Merkle root of the results from executing prev block
	LastResultsHash []byte `protobuf:"bytes,12,opt,name=last_results_hash,json=lastResultsHash,proto3" json:"last_results_hash,omitempty"`
	// the latest AppHash we've received from calling abci.Commit()
	AppHash              []byte   `protobuf:"bytes,13,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{4}
}
func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetVersion() Version {
	if m != nil {
		return m.Version
	}
	return Version{}
}

func (m *State) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *State) GetLastBlockHeight() int64 {
	if m != nil {
		return m.LastBlockHeight
	}
	return 0
}

func (m *State) GetLastBlockID() types1.BlockID {
	if m != nil {
		return m.LastBlockID
	}
	return types1.BlockID{}
}

func (m *State) GetLastBlockTime() time.Time {
	if m != nil {
		return m.LastBlockTime
	}
	return time.Time{}
}

func (m *State) GetNextValidators() *types1.ValidatorSet {
	if m != nil {
		return m.NextValidators
	}
	return nil
}

func (m *State) GetValidators() *types1.ValidatorSet {
	if m != nil {
		return m.Validators
	}
	return nil
}

func (m *State) GetLastValidators() *types1.ValidatorSet {
	if m != nil {
		return m.LastValidators
	}
	return nil
}

func (m *State) GetLastHeightValidatorsChanged() int64 {
	if m != nil {
		return m.LastHeightValidatorsChanged
	}
	return 0
}

func (m *State) GetConsensusParams() types1.ConsensusParams {
	if m != nil {
		return m.ConsensusParams
	}
	return types1.ConsensusParams{}
}

func (m *State) GetLastHeightConsensusParamsChanged() int64 {
	if m != nil {
		return m.LastHeightConsensusParamsChanged
	}
	return 0
}

func (m *State) GetLastResultsHash() []byte {
	if m != nil {
		return m.LastResultsHash
	}
	return nil
}

func (m *State) GetAppHash() []byte {
	if m != nil {
		return m.AppHash
	}
	return nil
}

func init() {
	proto.RegisterType((*ABCIResponses)(nil), "tendermint.proto.state.ABCIResponses")
	proto.RegisterType((*ValidatorsInfo)(nil), "tendermint.proto.state.ValidatorsInfo")
	proto.RegisterType((*ConsensusParamsInfo)(nil), "tendermint.proto.state.ConsensusParamsInfo")
	proto.RegisterType((*Version)(nil), "tendermint.proto.state.Version")
	proto.RegisterType((*State)(nil), "tendermint.proto.state.State")
}

func init() { proto.RegisterFile("proto/state/types.proto", fileDescriptor_00e69fef8162ea9b) }

var fileDescriptor_00e69fef8162ea9b = []byte{
	// 729 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x5d, 0x6a, 0xdb, 0x4a,
	0x18, 0xbd, 0xba, 0x4e, 0x62, 0xfb, 0x53, 0x1c, 0xdf, 0x3b, 0x81, 0x5c, 0x5d, 0x07, 0x6a, 0xe3,
	0x86, 0xc4, 0x2d, 0x45, 0x86, 0x74, 0x01, 0xa5, 0xb2, 0x4b, 0xa3, 0x92, 0x96, 0xa2, 0x84, 0x10,
	0xfa, 0x22, 0xc6, 0xd6, 0x44, 0x12, 0xb5, 0x25, 0xa1, 0x19, 0xbb, 0xc9, 0x1a, 0xfa, 0xd2, 0x1d,
	0x74, 0x3b, 0x5d, 0x85, 0x0b, 0x79, 0xee, 0x22, 0xca, 0xfc, 0x48, 0x9e, 0xfc, 0x11, 0x0c, 0x7d,
	0xf2, 0x68, 0xce, 0x77, 0xce, 0x77, 0x66, 0xe6, 0x7c, 0x18, 0xfe, 0xcb, 0xf2, 0x94, 0xa5, 0x7d,
	0xca, 0x30, 0x23, 0x7d, 0x76, 0x95, 0x11, 0x6a, 0x8b, 0x1d, 0xb4, 0xc3, 0x48, 0x12, 0x90, 0x7c,
	0x1a, 0x27, 0x4c, 0xee, 0xd8, 0xa2, 0xa6, 0xb5, 0xcf, 0xa2, 0x38, 0x0f, 0xfc, 0x0c, 0xe7, 0xec,
	0xaa, 0x2f, 0xc9, 0x61, 0x1a, 0xa6, 0xcb, 0x95, 0xac, 0x6e, 0xed, 0xe0, 0xd1, 0x38, 0x96, 0x8a,
	0xba, 0x6e, 0x4b, 0x35, 0xbc, 0x0b, 0xec, 0xea, 0xc0, 0x1c, 0x4f, 0xe2, 0x00, 0xb3, 0x34, 0x57,
	0xa0, 0xa5, 0x83, 0x19, 0xce, 0xf1, 0xf4, 0x16, 0x6d, 0x4e, 0x72, 0x1a, 0xa7, 0x49, 0xf1, 0xab,
	0xc0, 0x76, 0x98, 0xa6, 0xe1, 0x84, 0x48, 0x9f, 0xa3, 0xd9, 0x45, 0x9f, 0xc5, 0x53, 0x42, 0x19,
	0x9e, 0x66, 0xb2, 0xa0, 0xfb, 0xcb, 0x80, 0xc6, 0x6b, 0x67, 0xe0, 0x7a, 0x84, 0x66, 0x69, 0x42,
	0x09, 0x45, 0x2e, 0x98, 0x01, 0x99, 0xc4, 0x73, 0x92, 0xfb, 0xec, 0x92, 0x5a, 0x46, 0xa7, 0xd2,
	0x33, 0x0f, 0x7b, 0xb6, 0x76, 0x1b, 0xfc, 0x60, 0xb6, 0x74, 0x5e, 0xd0, 0x86, 0x92, 0x71, 0x7a,
	0xe9, 0x41, 0x50, 0x2c, 0x29, 0x1a, 0x42, 0x9d, 0x24, 0x81, 0x3f, 0x9a, 0xa4, 0xe3, 0xcf, 0xd6,
	0xdf, 0x1d, 0xa3, 0x67, 0x1e, 0x1e, 0x3c, 0x22, 0xf4, 0x26, 0x09, 0x1c, 0x5e, 0xee, 0xd5, 0x88,
	0x5a, 0xa1, 0x77, 0x60, 0x8e, 0x48, 0x18, 0x27, 0x4a, 0xa7, 0x22, 0x74, 0x9e, 0x3d, 0xa2, 0xe3,
	0x70, 0x86, 0x54, 0x82, 0x51, 0xb9, 0xee, 0x7e, 0x35, 0x60, 0xeb, 0xac, 0xb8, 0x5a, 0xea, 0x26,
	0x17, 0x29, 0x72, 0xa1, 0x51, 0x5e, 0xb6, 0x4f, 0x09, 0xb3, 0x0c, 0xd1, 0x60, 0xcf, 0xbe, 0xf3,
	0xfe, 0xb2, 0x43, 0x49, 0x3f, 0x21, 0xcc, 0xdb, 0x9c, 0x6b, 0x5f, 0xc8, 0x86, 0xed, 0x09, 0xa6,
	0xcc, 0x8f, 0x48, 0x1c, 0x46, 0xcc, 0x1f, 0x47, 0x38, 0x09, 0x49, 0x20, 0x4e, 0x5e, 0xf1, 0xfe,
	0xe5, 0xd0, 0x91, 0x40, 0x06, 0x12, 0xe8, 0x7e, 0x37, 0x60, 0x7b, 0xc0, 0xdd, 0x26, 0x74, 0x46,
	0x3f, 0x8a, 0x47, 0x15, 0x96, 0xce, 0xe1, 0x9f, 0x71, 0xb1, 0xed, 0xcb, 0xc7, 0x56, 0xae, 0x0e,
	0x1e, 0x72, 0x75, 0x4b, 0xc6, 0x59, 0xfb, 0xb1, 0x68, 0xff, 0xe5, 0x35, 0xc7, 0x37, 0xb7, 0x57,
	0x76, 0x98, 0x40, 0xf5, 0x4c, 0x06, 0x0a, 0xbd, 0x85, 0x7a, 0xa9, 0xa6, 0xdc, 0x3c, 0xbd, 0xeb,
	0xa6, 0x88, 0x5f, 0xe9, 0x47, 0x39, 0x59, 0x72, 0x51, 0x0b, 0x6a, 0x34, 0xbd, 0x60, 0x5f, 0x70,
	0x4e, 0x44, 0xe3, 0xba, 0x57, 0x7e, 0x77, 0x17, 0x1b, 0xb0, 0x7e, 0xc2, 0xc7, 0x0c, 0xbd, 0x82,
	0xaa, 0xd2, 0x52, 0xcd, 0xda, 0xf6, 0xfd, 0x03, 0x69, 0x2b, 0x83, 0xaa, 0x51, 0xc1, 0x42, 0xfb,
	0x50, 0x1b, 0x47, 0x38, 0x4e, 0xfc, 0x58, 0x9e, 0xaf, 0xee, 0x98, 0xd7, 0x8b, 0x76, 0x75, 0xc0,
	0xf7, 0xdc, 0xa1, 0x57, 0x15, 0xa0, 0x1b, 0xa0, 0xe7, 0x20, 0xce, 0x2d, 0xd3, 0xa5, 0x2e, 0x46,
	0x84, 0xac, 0xe2, 0x35, 0x39, 0x20, 0x82, 0x23, 0x6f, 0x05, 0x9d, 0x43, 0x43, 0xab, 0x8d, 0x03,
	0x6b, 0xed, 0x21, 0x6b, 0xf2, 0x55, 0x04, 0xd7, 0x1d, 0x3a, 0xdb, 0xdc, 0xda, 0xf5, 0xa2, 0x6d,
	0x1e, 0x17, 0x82, 0xee, 0xd0, 0x33, 0x4b, 0x75, 0x37, 0x40, 0xc7, 0xd0, 0xd4, 0x94, 0xf9, 0x94,
	0x5a, 0xeb, 0x42, 0xbb, 0x65, 0xcb, 0x11, 0xb6, 0x8b, 0x11, 0xb6, 0x4f, 0x8b, 0x11, 0x76, 0x6a,
	0x5c, 0xf6, 0xdb, 0xcf, 0xb6, 0xe1, 0x35, 0x4a, 0x2d, 0x8e, 0xa2, 0xf7, 0xd0, 0x4c, 0xc8, 0x25,
	0xf3, 0xcb, 0x74, 0x52, 0x6b, 0x63, 0x85, 0x54, 0x6f, 0x71, 0xf2, 0x72, 0x4c, 0xd0, 0x10, 0x40,
	0x53, 0xaa, 0xae, 0xa0, 0xa4, 0xf1, 0xb8, 0x29, 0x71, 0x44, 0x4d, 0xaa, 0xb6, 0x8a, 0x29, 0x4e,
	0xd6, 0x4c, 0x0d, 0xe0, 0x89, 0x1e, 0xe5, 0xa5, 0x6a, 0x99, 0xea, 0xba, 0x78, 0xc4, 0xdd, 0x65,
	0xaa, 0x97, 0x6c, 0x95, 0xef, 0x7b, 0x27, 0x0d, 0xfe, 0xc8, 0xa4, 0x7d, 0x80, 0xbd, 0x1b, 0x93,
	0x76, 0xab, 0x4b, 0x69, 0xd2, 0x14, 0x26, 0x3b, 0xda, 0xe8, 0xdd, 0x14, 0x2a, 0x9c, 0x16, 0x31,
	0xcd, 0x09, 0x9d, 0x4d, 0x18, 0xf5, 0x23, 0x4c, 0x23, 0x6b, 0xb3, 0x63, 0xf4, 0x36, 0x65, 0x4c,
	0x3d, 0xb9, 0x7f, 0x84, 0x69, 0x84, 0xfe, 0x87, 0x1a, 0xce, 0x32, 0x59, 0xd2, 0x10, 0x25, 0x55,
	0x9c, 0x65, 0x1c, 0x72, 0xec, 0x4f, 0x2f, 0xc2, 0x98, 0x45, 0xb3, 0x91, 0x3d, 0x4e, 0xa7, 0xfd,
	0xe5, 0x11, 0xf5, 0xa5, 0xf6, 0x8f, 0x38, 0xda, 0x10, 0x1f, 0x2f, 0x7f, 0x07, 0x00, 0x00, 0xff,
	0xff, 0x93, 0x33, 0x0f, 0xa0, 0x27, 0x07, 0x00, 0x00,
}
