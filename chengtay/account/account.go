package account

import (
	"bytes"
	"errors"
	"time"

	"github.com/tendermint/tendermint/crypto"
)

//-----------------------------------------------------------------------------
type Coin struct {
	Denom string `json:"denom"`

	// To allow the use of unsigned integers (see: #1273) a larger refactor will
	// need to be made. So we use signed integers for now with safety measures in
	// place preventing negative values being used.
	Amount Int `json:"amount"`
}
type Coins []Coin
// implements Account.
type BaseAccount struct {
	Address       AccAddress    `json:"address" yaml:"address"`
	Coins         Coins         `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey `json:"public_key" yaml:"public_key"`
	AccountNumber uint64        `json:"account_number" yaml:"account_number"`
	Sequence      uint64        `json:"sequence" yaml:"sequence"`
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address AccAddress, coins Coins,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}


func (acc BaseAccount) GetAddress() AccAddress {
	return acc.Address
}


func (acc *BaseAccount) SetAddress(addr AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

func (acc *BaseAccount) GetCoins() Coins {
	return acc.Coins
}

func (acc *BaseAccount) SetCoins(coins Coins) error {
	acc.Coins = coins
	return nil
}

func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}


func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) Coins {
	return acc.GetCoins()
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && acc.Address != nil &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}
