package test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
)

//生成随即字符串
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func Genbyte(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}
func GenString(n int) string {
	b := make([]byte, n)
	b = Genbyte(n)
	return *(*string)(unsafe.Pointer(&b))
}
func GenInt() uint64 {
	rand.Seed(666)
	return rand.Uint64()
}

//merkle tree node都填满，内容随机，生成hash，对每个节点生成proof，验证proof
//测试情景二，部分节点不填满，应该会在生成hash时报错
//测试情景三，一部分节点随机，另一部分节点用DummyMerkleNode填充，对每个节点生成proof，验证proof

func TestMerkletree1(t *testing.T) {
	StorageItems := [4]types.StorageItem{}
	for i := 0; i < 4; i++ {
		StorageItems[i].CarID = (types.ID(GenString(16)))
		StorageItems[i].Content = Genbyte(16)
		StorageItems[i].StorageItemID = (types.ID(GenString(16)))
		StorageItems[i].Timestamp = GenInt()
	}
	var StorageItemMerkleNode [4]types.StorageItemMerkleNode
	for i := 0; i < 4; i++ {
		StorageItemMerkleNode[i] = types.StorageItemMerkleNode(StorageItems[i])
	}

	//默克尔树插入节点
	MerkleTree := new(types.MerkleTree)
	for i := 0; i < 4; i++ {
		err := MerkleTree.SetMerkleNode(i, &StorageItemMerkleNode[i])
		if err != nil {
			panic(err)
		}
	}
	//测试
	fmt.Println(MerkleTree.GetCapacity())
	//所有的proof应该是相等的
	var MerkleProof [4]types.MerkleProof
	var err error
	for i := 0; i < 4; i++ {
		MerkleProof[i], err = MerkleTree.GetNodeMerkleProof(i)
		if err != nil {
			panic(err)
		}
		//fmt.Println(MerkleProof[i].MerkleRoot)
	}
	a := bytes.Equal(MerkleProof[1].MerkleRoot, MerkleProof[2].MerkleRoot)
	b := bytes.Equal(MerkleProof[2].MerkleRoot, MerkleProof[3].MerkleRoot)
	c := bytes.Equal(MerkleProof[3].MerkleRoot, MerkleProof[0].MerkleRoot)
	fmt.Println(a, b, c)

	//只设置3个节点，不full的情况，//测试情景二，部分节点不填满，应该会在生成hash时报错//
	MerkleTree2 := new(types.MerkleTree)
	for i := 0; i < 3; i++ {
		err := MerkleTree2.SetMerkleNode(i, &StorageItemMerkleNode[i])
		if err != nil {
			panic(err)
		}
	}
	//测试
	var MerkleProof1 types.MerkleProof
	MerkleProof1, err = MerkleTree2.GetNodeMerkleProof(0)
	//输出Merkle tree is not full. Can't compute merkle root.
	fmt.Println(err, MerkleProof1)

	/*
	//测试情景三，一部分节点随机，另一部分节点用DummyMerkleNode填充，对每个节点生成proof，验证proof
	//节点5  DummyMerkleNode节点
	var DummyMerkleNode types.IMerkleNode = new(types.DummyMerkleNode)
	fmt.Println(DummyMerkleNode)
	err = MerkleTree2.SetMerkleNode(3, DummyMerkleNode)
	fmt.Println(MerkleProof1.Validate())
	 */
}
