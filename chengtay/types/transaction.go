package types

import (
	"bytes"
	"encoding/json"
	"github.com/ChengtayChain/ChengtayChain/chengtay/util"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	"math/big"
)

type RawTransaction struct {
	Type               uint32
	Value              []byte                // json string that contains a ITransactionValue
	ValueHash          []byte                // sha256 of Value
	ValueHashSignature []byte                // the signature of ValueHash
	PublicKey          ed25519.PubKeyEd25519 // the corresponding public key of the signature
}

func (self *RawTransaction) VerifyHash() bool {
	return bytes.Equal(DefaultHashProvider.Digest(self.Value), self.ValueHash)
}

func (self *RawTransaction) VerifySignature() bool {
	return self.VerifyHash() && self.PublicKey.VerifyBytes(self.ValueHash, self.ValueHashSignature)
}

type Transaction struct {
	Type      uint32
	PublicKey ed25519.PubKeyEd25519
	Value     ITransactionValue
}

type ITransactionValue interface {
	GetType() uint32
	GetTimestamp() uint64 // Unix timestamp in second
	GetHash() (digest []byte, err error)
}

const (
	TransactionUnknown                = 0
	TransactionMerkleroot             = 1
	TransactionUpdateTrustedPublicKey = 2
	//TransactionVehicleUploading
	TransactionRequestVehicle       = 4
	TransactionResponseVehicle      = 5
	TransactionResponseTimeout      = 6
	TransactionReturnVehicleTimeout = 7
	TransactionClearing             = 8
	TransactionProof                = 9
)

func GetTransaction(raw RawTransaction) (errorCode uint32, ret Transaction, criticalError error) {
	// TODO: use a better method to replace if-else statements?
	if raw.Type == TransactionMerkleroot {
		var value MerkleRootTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionUpdateTrustedPublicKey {
		// TODO: not implemented yet
		return ErrorNotImplemented, Transaction{}, nil
	} else if raw.Type == TransactionRequestVehicle {
		var value RequestVehicleTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionResponseVehicle {
		var value ResponseVehicleTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionResponseTimeout {
		var value ResponseTimeoutTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionReturnVehicleTimeout {
		var value ReturnVehicleTimeoutTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionClearing {
		var value ClearingTransactionValue
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else if raw.Type == TransactionProof {
		var value ZkProofTransaction
		err := json.Unmarshal(raw.Value, &value)
		if err != nil {
			return ErrorJsonParsing, Transaction{}, nil
		}
		ret.Type = raw.Type
		ret.PublicKey = raw.PublicKey
		ret.Value = &value
		return ErrorNoError, ret, nil
	} else {
		return ErrorUnknownTransactionType, Transaction{}, nil
	}
}

//MerkleRootTransactionValue
type MerkleRootTransactionValue struct {
	Timestamp uint64
	Nonce     [256]byte // random bytes

	Items []MerkleRootTransactionItem
}

type MerkleRootTransactionItem []byte

func (self *MerkleRootTransactionValue) GetType() uint32 {
	return TransactionMerkleroot
}

func (self *MerkleRootTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp
}

func (self *MerkleRootTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}

	// 3. MerkleRootTransactionItem([]byte) in order
	for i := 0; i < len(self.Items); i++ {
		source = append(source, self.Items[i]...)
	}

	return DefaultHashProvider.Digest(source), nil
}

//RequestVehicleTransactionValue type3
type RequestVehicleTransactionValue struct {
	Timestamp   uint64
	Nonce       [256]byte // random bytes
	VehicleID   string
	PublicKey   ed25519.PubKeyEd25519
	Deposit     *big.Int //
	RequestTime string   //compute by day
	//zk
	cipherDeposit []byte
	ZKProof       []byte //cipherDeposit  == encrypt(Deposit)
}

func (self *RequestVehicleTransactionValue) GetType() uint32 {
	return TransactionRequestVehicle
}

func (self *RequestVehicleTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp
}

func (self *RequestVehicleTransactionValue) GetVehicleID() string {
	return self.VehicleID
}

func (self *RequestVehicleTransactionValue) GetPublicKey() ed25519.PubKeyEd25519 {
	return self.PublicKey
}
func (self *RequestVehicleTransactionValue) GetRent() *big.Int {
	return self.Deposit
}
func (self *RequestVehicleTransactionValue) GetRequestTime() string {
	return self.RequestTime
}
func (self *RequestVehicleTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}

	//3.VehicleID
	source = append(source, self.VehicleID...)

	//4.PublicKey
	source = append(source, self.PublicKey[:]...)
	//5.Deposit
	source = append(source, self.Deposit.Bytes()...)

	//6.RequestTime
	source = append(source, self.RequestTime...)
	return DefaultHashProvider.Digest(source), nil
}

//ResponseVehicleTransactionValue type4
type ResponseVehicleTransactionValue struct {
	Timestamp          uint64
	Nonce              [256]byte // random bytes
	VehicleID          string
	VehicleIDSignature []byte
	PublicKey          ed25519.PubKeyEd25519 // the corresponding public key of the signature
}

func (self *ResponseVehicleTransactionValue) VerifySignature() bool {
	return self.PublicKey.VerifyBytes([]byte(self.VehicleID), self.VehicleIDSignature) //if true ,means unlock the car
}
func (self *ResponseVehicleTransactionValue) GetPublicKey() ed25519.PubKeyEd25519 {
	return self.PublicKey
}
func (self *ResponseVehicleTransactionValue) GetType() uint32 {
	return TransactionResponseVehicle
}

func (self *ResponseVehicleTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp
}
func (self *ResponseVehicleTransactionValue) GetVehicleID() string {
	return self.VehicleID
}
func (self *ResponseVehicleTransactionValue) GetVehicleIDSignature() []byte {
	return self.VehicleIDSignature
}

func (self *ResponseVehicleTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}
	//3.VehicleID
	source = append(source, self.VehicleID...)
	//4.VehicleIDSignature
	source = append(source, self.VehicleIDSignature...)
	//5.PublicKey
	source = append(source, self.PublicKey[:]...)
	return DefaultHashProvider.Digest(source), nil
}

//ResponseTimeoutTransactionValue  type5
//TODO uncomplete how to handle time?
type ResponseTimeoutTransactionValue struct {
	Timestamp uint64
	Nonce     [256]byte // random bytes
	VehicleID string
	//ExceededTime uint64

	PublicKey ed25519.PubKeyEd25519 // the corresponding public key of the signature
}

func (self *ResponseTimeoutTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp //if true ,means unlock the car
}
func (self *ResponseTimeoutTransactionValue) GetType() uint32 {
	return TransactionResponseTimeout
}
func (self *ResponseTimeoutTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}
	//3.VehicleID
	source = append(source, self.VehicleID[:]...)
	//5.PublicKey
	source = append(source, self.PublicKey[:]...)
	return DefaultHashProvider.Digest(source), nil
}

//ReturnVehicleTimeoutTransactionValue  type6 TransactionReturnVehicleTimeout
//TODO uncomplete how to handle time?
type ReturnVehicleTimeoutTransactionValue struct {
	Timestamp    uint64
	Nonce        [256]byte // random bytes
	VehicleID    string
	ExceededTime uint64
	PublicKey    ed25519.PubKeyEd25519 // the corresponding public key of the signature to verify
	//UserPublicKey   ed25519.PubKeyEd25519 //Return Timeout so add to the untrusted list
}

func (self *ReturnVehicleTimeoutTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp //if true ,means unlock the car
}
func (self *ReturnVehicleTimeoutTransactionValue) GetType() uint32 {
	return TransactionReturnVehicleTimeout
}
func (self *ReturnVehicleTimeoutTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}
	//3.VehicleID
	source = append(source, self.VehicleID[:]...)
	//5.PublicKey
	source = append(source, self.PublicKey[:]...)
	return DefaultHashProvider.Digest(source), nil
}

//ClearingTransactionValue type7
type ClearingTransactionValue struct {
	Timestamp uint64
	Nonce     [256]byte // random bytes
	VehicleID string

	//UsedTime  uint64
	//MultiSig
	//PublicKey          ed25519.PubKeyEd25519 // the corresponding public key of the signature

	///zk
	ZKProof []byte
}

func (self *ClearingTransactionValue) GetTimestamp() uint64 {
	return self.Timestamp //if true ,means unlock the car
}
func (self *ClearingTransactionValue) GetType() uint32 {
	return TransactionClearing
}
func (self *ClearingTransactionValue) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}
	//3.VehicleID
	source = append(source, self.VehicleID[:]...)

	return DefaultHashProvider.Digest(source), nil
}

//type ZkProofTransaction struct {
//	Proof   groth16.Proof
//	Vk      groth16.VerifyingKey
//	Witness zktx.PaillerCircuit
//}

//////////////
//use xjsnark
type ZkProofTransaction struct {
	Timestamp uint64
	Nonce     [256]byte // random bytes
	Proof     []byte
}

func (self *ZkProofTransaction) GetType() uint32 {
	return TransactionProof
}

func (self *ZkProofTransaction) GetTimestamp() uint64 {
	return self.Timestamp
}

func (self *ZkProofTransaction) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. Nonce
	source = append(source, self.Nonce[:]...)

	// 2. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}

	return DefaultHashProvider.Digest(source), nil
}
