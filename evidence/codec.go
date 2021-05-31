package evidence

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/ChengtayChain/ChengtayChain/crypto/encoding/amino"
	"github.com/ChengtayChain/ChengtayChain/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
	cryptoamino.RegisterAmino(cdc)
	types.RegisterEvidences(cdc)
}

// For testing purposes only
func RegisterMockEvidences() {
	types.RegisterMockEvidences(cdc)
}
