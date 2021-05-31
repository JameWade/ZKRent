package consensus

import (
	amino "github.com/tendermint/go-amino"

	"github.com/ChengtayChain/ChengtayChain/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
	RegisterWALMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
