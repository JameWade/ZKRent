package client

import (
	amino "github.com/tendermint/go-amino"

	"github.com/ChengtayChain/ChengtayChain/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterEvidences(cdc)
}
