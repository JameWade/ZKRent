package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	tmos "github.com/ChengtayChain/ChengtayChain/libs/os"
	"github.com/ChengtayChain/ChengtayChain/libs/rand"
	"github.com/ChengtayChain/ChengtayChain/privval"
	ctypes "github.com/ChengtayChain/ChengtayChain/rpc/core/types"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"os"
	"time"

	rpcClient "github.com/ChengtayChain/ChengtayChain/rpc/client/http"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	privKey := getPriKey()
	var address = "http://127.0.0.1:26657" // cfg.DefaultRPCConfig().ListenAddress
	var client *rpcClient.HTTP
	{
		var err error
		client, err = rpcClient.New(address, "/websocket")
		if err != nil {
			panic(err)
		}
	}

	for {
		rawTransaction := randomRawTransaction(privKey)
		ret, err := sendRawTransaction(client, rawTransaction)
		if err != nil {
			fmt.Println(err.Error())
			// ignore the error
		} else {
			fmt.Printf("%+v\n", ret)
		}

		time.Sleep(15 * time.Second)
	}

}

func sendRawTransaction(client *rpcClient.HTTP, rawTransaction types.RawTransaction) (*ctypes.ResultBroadcastTx, error) {
	bytes, err := json.Marshal(&rawTransaction)
	if err != nil {
		panic(err)
	}

	ret, err := client.BroadcastTxSync(bytes)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func randomRawTransaction(privKey ed25519.PrivKeyEd25519) (rawTransaction types.RawTransaction) {
	filepath := "/home/waris/circuit/proof"
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var proof = content
	value := types.ZkProofTransaction{
		Timestamp: uint64(time.Now().Unix()),
		Proof:     proof,
	}
	bytes256 := rand.Bytes(256)
	copy(value.Nonce[:256], bytes256[:256])
	////////////////feihua
	valueBytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	valueBytesHash := types.DefaultHashProvider.Digest(valueBytes)

	valueBytesHashSig, err := privKey.Sign(valueBytesHash)
	if err != nil {
		panic(err)
	}

	rawTransaction.Type = types.TransactionProof
	rawTransaction.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	rawTransaction.Value = valueBytes
	rawTransaction.ValueHash = valueBytesHash
	rawTransaction.ValueHashSignature = valueBytesHashSig

	return rawTransaction
}
func getPriKey() ed25519.PrivKeyEd25519 {
	// get TMHOME
	var tmhome string
	{
		tmhome = os.Getenv("TMHOME")
		if len(tmhome) == 0 {
			homeDir, err := homedir.Dir()
			if err != nil {
				panic(err)
			}
			tmhome = homeDir + string(os.PathSeparator) + ".chengtaychain"
		}
	}

	// get private key
	var privKey ed25519.PrivKeyEd25519
	{
		privValiKeyFile := tmhome + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "priv_validator_key.json"
		privValiStateFile := tmhome + string(os.PathSeparator) + "data" + string(os.PathSeparator) + "priv_validator_state.json"
		var pv *privval.FilePV
		if tmos.FileExists(privValiKeyFile) {
			pv = privval.LoadFilePV(privValiKeyFile, privValiStateFile)
		} else {
			panic(fmt.Errorf("file not found. " + privValiKeyFile))
		}

		privKey = pv.Key.PrivKey.(ed25519.PrivKeyEd25519)
	}
	return privKey
}
