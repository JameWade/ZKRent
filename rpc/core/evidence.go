package core

import (
	"github.com/ChengtayChain/ChengtayChain/evidence"
	ctypes "github.com/ChengtayChain/ChengtayChain/rpc/core/types"
	rpctypes "github.com/ChengtayChain/ChengtayChain/rpc/jsonrpc/types"
	"github.com/ChengtayChain/ChengtayChain/types"
)

// BroadcastEvidence broadcasts evidence of the misbehavior.
// More: https://docs.tendermint.com/master/rpc/#/Info/broadcast_evidence
func BroadcastEvidence(ctx *rpctypes.Context, ev types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	err := env.EvidencePool.AddEvidence(ev)
	if _, ok := err.(evidence.ErrEvidenceAlreadyStored); err == nil || ok {
		return &ctypes.ResultBroadcastEvidence{Hash: ev.Hash()}, nil
	}
	return nil, err
}
