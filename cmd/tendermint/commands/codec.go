package commands

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/ChengtayChain/ChengtayChain/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
}
