package test

import (
	"encoding/json"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	cfg "github.com/ChengtayChain/ChengtayChain/config"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	tmos "github.com/ChengtayChain/ChengtayChain/libs/os"
	"github.com/ChengtayChain/ChengtayChain/privval"
	"io/ioutil"
	"testing"

	rpcclient "github.com/ChengtayChain/ChengtayChain/rpc/client/http"
)


// KEYFILENAME 私钥文件名
const PRIVATEKEYFILENAME =     "./private.key.ed25519.json"
const PUBLICKEYFILENAME =     "./public.key.ed25519.json"
var (
	cli *rpcclient.HTTP
	me  *user
)

func init() {
	var err error
	addr := cfg.DefaultRPCConfig().ListenAddress
	cli, err = rpcclient.New(addr, "/websocket")
	if err != nil {
		return
	}
	user, err := loadUserKey()
	if err != nil {
		panic(err)
	}
	me = user
}
  /*
type PubKey struct {
	types string        `json:type`
	value []byte     `json:value`
}
type PrivKey struct {
	types string       `json:type`
	value []byte        `json:value`
}
type validator struct{
	address string`json:address`
	pubKey   PubKey     `json:value`
	privkey    PrivKey     `json:privkey`
}
  */

type user struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  ed25519.PubKeyEd25519
}

func loadUserKey() (*user, error) {
	/*copy(privKey[:], bz)
	jsonBytes, err := ioutil.ReadFile(PRIVATEKEYFILENAME)
	if err != nil {
		return nil, err
	}
	//var v validator
	fmt.Println(jsonBytes)
	err = json.Unmarshal(jsonBytes, &v)
	if err != nil {
		panic(err)
	}
	//fmt.Println(v.address)
	//uk.PubKey =    uk.PrivKey.PubKey()
*/
	 privValKeyFile :=   "/home/wjp/.chengtaychain/config/priv_validator_key.json"
	 privValStateFile :=    "/home/wjp/.chengtaychain/data/priv_validator_state.json"
	 var pv *privval.FilePV
	 if tmos.FileExists(privValKeyFile) {
	     pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
	  }
	uk := new(user)
	/*
	var bytes32 [32]byte
	 for k, v := range pv.Key.PubKey.Bytes() {
	 	bytes32[k] = v
	 }
	uk.PubKey =  ed25519.PubKeyEd25519(bytes32 )
	var bytes64 [ed25519.SignatureSize]byte
	 for k, v := range pv.Key.PrivKey.Bytes() {
	 	bytes64[k] = v
	 }
	uk.PrivKey = ed25519.PrivKeyEd25519(bytes64 )

	 */
	  {
	  		bytes, err := json.Marshal(pv.Key.PrivKey)
	  		if err != nil {
	  			panic(err)
	  		}
	  		err = ioutil.WriteFile("private.key.ed25519.json", bytes, 0644)
	  		if err != nil {
	  			panic(err)
	  		}
	  	}
	  	{
	  		bytes, err := json.Marshal(pv.Key.PubKey)
	  		if err != nil {
	  			panic(err)
	  		}
	  		err = ioutil.WriteFile("public.key.ed25519.json", bytes, 0644)
	  		if err != nil {
	  			panic(err)
	  		}
	  	}
	jsonBytes, err := ioutil.ReadFile(PRIVATEKEYFILENAME)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &uk.PrivKey)
		jsonBytes, err = ioutil.ReadFile(PUBLICKEYFILENAME)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &uk.PubKey)
		if err != nil {
			panic(err)
		}

	fmt.Println(uk.PrivKey)
	fmt.Println(uk.PubKey)
	fmt.Println(uk.PrivKey.PubKey() )
	return uk, nil
}
func GetMerkleRoot(index int) ([]byte){
	StorageItems := [4](types.StorageItemMerkleNode){}
	for i := 0; i < 4; i++ {
		StorageItems[i].CarID = (types.ID(GenString(16)))
		StorageItems[i].Content = Genbyte(16)
		StorageItems[i].StorageItemID = (types.ID(GenString(16)))
		StorageItems[i].Timestamp = GenInt()
	}
	var StorageItemMerkleNode [4]types.IMerkleNode
	for i := 0; i < 4; i++ {
		StorageItemMerkleNode[i] = &StorageItems[i]
	}

	//默克尔树插入节点
	MerkleTree := new(types.MerkleTree)
	for i := 0; i < 4; i++ {
		err := MerkleTree.SetMerkleNode(i, StorageItemMerkleNode[i])
		if err != nil {
			panic(err)
		}
	}

	var MerkleProof [4]types.MerkleProof
	var err error
	for i := 0; i < 4; i++ {
		MerkleProof[i], err = MerkleTree.GetNodeMerkleProof(i)
		if err != nil {
			panic(err)
		}
		//fmt.Println(MerkleProof[i].MerkleRoot)
	}
	return MerkleProof[index].MerkleRoot
}


func (bg *user) genTx() error {
	//交易填充
	var bytes []byte
	var ValueSignature []byte
	var MerkleRootTransactionItem types.MerkleRootTransactionItem
	var Items [4]types.MerkleRootTransactionItem
	MerkleRootTransactionValue := new(types.MerkleRootTransactionValue)
	MerkleRootTransactionValue.Nonce = [256]byte{123}
	MerkleRootTransactionValue.Timestamp = 123456      //why not use type time
	MerkleRootTransactionItem = GetMerkleRoot(0)
	Items[0] = append(Items[0], MerkleRootTransactionItem...)
	MerkleRootTransactionValue.Items = append(MerkleRootTransactionValue.Items,Items[0])
	valueBytes, err := json.Marshal(MerkleRootTransactionValue)   //only one
	valueBytesHash := types.DefaultHashProvider.Digest(valueBytes)
	tx := new(types.RawTransaction)
	tx.PublicKey = bg.PubKey
	tx.Value = bytes
	//交易签名
	ValueSignature,err = bg.PrivKey.Sign(valueBytesHash)
	if err!=nil{
		panic(err)
	}
	tx.ValueHash = valueBytesHash
	tx.ValueHashSignature = ValueSignature
	//交易广播
	bz, err := json.Marshal(&tx)
	fmt.Println("verify signature before send ",tx.VerifySignature())
	if err != nil {
		return err
	}
	ret, err := cli.BroadcastTxSync(bz)
	if err != nil {
		return err
	}

	fmt.Printf("throw A Tx=> %+v\n", ret)

	return nil
}
func TestTransaction(t *testing.T)  {
	for {
	err := me.genTx()
	if err != nil {
		fmt.Print(err)
	}
		}
}

