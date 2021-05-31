package types

import (
	"fmt"
)

type IMerkleTree interface {
	GetMerkleNode(index int) (node IMerkleNode, err error)
	SetMerkleNode(index int, node IMerkleNode) (err error)
	GetMerkleRoot() (root []byte, err error)
	GetNodeMerkleProof(index int) (proof MerkleProof, err error)
	GetCapacity() int
}

const merkleTreeCapacity int = 4 // must be power of 2

type MerkleTree struct {
	nodes [merkleTreeCapacity]IMerkleNode
}

func (self *MerkleTree) GetCapacity() int {
	return merkleTreeCapacity
}

func (self *MerkleTree) GetMerkleNode(index int) (node IMerkleNode, err error) {
	return self.nodes[index], nil
}

func (self *MerkleTree) SetMerkleNode(index int, node IMerkleNode) (err error) {
	self.nodes[index] = node
	return nil
}

func (self *MerkleTree) GetMerkleRoot() (root []byte, err error) {
	proof, err := self.GetNodeMerkleProof(0)
	if err != nil {
		return nil, err
	}
	return proof.MerkleRoot, nil
}

func (self *MerkleTree) GetNodeMerkleProof(index int) (proof MerkleProof, err error) {
	proof.Node = self.nodes[index]
	proof.Path = make([]MerkleProofHash, 0)

	var hashes [merkleTreeCapacity][]byte
	for j := 0; j < merkleTreeCapacity; j++ {
		if self.nodes[j] == nil {
			return MerkleProof{}, fmt.Errorf("Merkle tree is not full. Can't compute merkle root.")
		}

		hashes[j], err = self.nodes[j].GetHash()
		if err != nil {
			return MerkleProof{}, err
		}
	}

	i := index
	n := merkleTreeCapacity

	for n >= 2 {
		var proofHash MerkleProofHash
		if i%2 == 0 {
			proofHash.IsLeft = true
			proofHash.Hash, err = self.nodes[i+1].GetHash()
			if err != nil {
				return MerkleProof{}, err
			}
		} else {
			proofHash.IsLeft = false
			proofHash.Hash, err = self.nodes[i-1].GetHash()
			if err != nil {
				return MerkleProof{}, err
			}
		}

		proof.Path = append(proof.Path, proofHash)

		for j := 0; j < n/2; j++ {
			hashes[j*2] = DefaultHashProvider.Digest(append(hashes[j*2], hashes[j*2+1]...))
		}

		n = n / 2

		i = i / 2
	}

	proof.MerkleRoot = hashes[0]

	return proof, nil
}
