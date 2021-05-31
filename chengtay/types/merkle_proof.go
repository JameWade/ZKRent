package types

import (
	"bytes"
)

type MerkleProofHash struct {
	Hash   []byte
	IsLeft bool
}

type MerkleProof struct {
	Node       IMerkleNode
	Path       []MerkleProofHash
	MerkleRoot []byte
}

func (self *MerkleProof) Validate() (ret bool, err error) {
	hash, err := self.Node.GetHash()
	if err != nil {
		return false, err
	}

	for i := 0; i < len(self.Path); i++ {
		if self.Path[i].IsLeft {
			hash = DefaultHashProvider.Digest(append(hash, self.Path[i].Hash...))
		} else {
			hash = DefaultHashProvider.Digest(append(self.Path[i].Hash, hash...))
		}
	}

	return bytes.Equal(self.MerkleRoot, hash), nil
}
