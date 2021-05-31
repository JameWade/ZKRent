package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	"github.com/ChengtayChain/ChengtayChain/libs/rand"
	rpcClient "github.com/ChengtayChain/ChengtayChain/rpc/client/http"
	ctypes "github.com/ChengtayChain/ChengtayChain/rpc/core/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

// sendtxCmd represents the sendtx command
var SendtxCmd = &cobra.Command{
	Use:   "sendtx",
	Short: "Send a  Tx",
	Long:  `Send a  Tx`,
	RunE:  sendTx,
}

func initialize() (privKey ed25519.PrivKeyEd25519) {
	//need account privkey
	// read public key file
	var err error
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var flagKeyDir string
	flag.StringVar(&flagKeyDir, "k", workingDir, "the directory which contains 'private.key.ed25519.json' file.")
	{
		bytes, err := ioutil.ReadFile(workingDir + string(os.PathSeparator) + "private.key.ed25519.json")
		if err != nil {
			fmt.Println("please run VehicleUser-client genkey firstly")
			panic(err)
		}
		err = json.Unmarshal(bytes, &privKey)

		if err != nil {
			panic(err)
		}
	}
	return privKey
}
func sendTx(cmd *cobra.Command, args []string) error {
	privKey := initialize()
	//send Tx
	var address = "http://127.0.0.1:26657" // cfg.DefaultRPCConfig().ListenAddress
	var client *rpcClient.HTTP
	{
		var err error
		client, err = rpcClient.New(address, "/websocket")
		if err != nil {
			panic(err)
		}
	}

	if args[0] == "ResquestVehicle" {
		for i := 0; i < 1; i++ {
			rawTransaction := randomRawTransaction(privKey, types.TransactionRequestVehicle)
			ret, err := sendRawTransaction(client, rawTransaction)
			if err != nil {
				fmt.Println(err.Error())
				// ignore the error
			} else {
				fmt.Printf("%+v\n", ret)
			}

			time.Sleep(500 * time.Millisecond)
		}
	} else if args[0] == "ResponseTimeout" {
		for i := 0; i < 1; i++ {
			rawTransaction := randomRawTransaction(privKey, types.TransactionResponseTimeout)
			ret, err := sendRawTransaction(client, rawTransaction)
			if err != nil {
				fmt.Println(err.Error())
				// ignore the error
			} else {
				fmt.Printf("%+v\n", ret)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
	return nil
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
func randomRequestVehicleTransactionValue(privKey ed25519.PrivKeyEd25519) (value *types.RequestVehicleTransactionValue) {
	value = &types.RequestVehicleTransactionValue{}
		value.Timestamp = uint64(time.Now().Unix())
	bytes256 := rand.Bytes(256)
	copy(value.Nonce[:256], bytes256[:256])
	value.VehicleID = "1"
	value.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	value.Deposit, _ = new(big.Int).SetString("200", 10)
	value.RequestTime = "1" //days
	return value
}

func randomResponseTimeoutTransactionValue(privKey ed25519.PrivKeyEd25519) (value *types.ResponseVehicleTransactionValue) {
	value = &types.ResponseVehicleTransactionValue{}
	value.Timestamp = uint64(time.Now().Unix())
	bytes256 := rand.Bytes(256)
	copy(value.Nonce[:256], bytes256[:256])
	value.VehicleID = "1"
	value.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	return value
}
func randomRawTransaction(privKey ed25519.PrivKeyEd25519, transactionType uint32) (rawTransaction types.RawTransaction) {
	var value types.ITransactionValue
	if transactionType == types.TransactionRequestVehicle {
		value = randomRequestVehicleTransactionValue(privKey)
	} else if transactionType == types.TransactionResponseTimeout {
		value = randomResponseTimeoutTransactionValue(privKey)
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	valueBytesHash := types.DefaultHashProvider.Digest(valueBytes)
	valueBytesHashSig, err := privKey.Sign(valueBytesHash)
	if err != nil {
		panic(err)
	}

	rawTransaction.Type = transactionType
	rawTransaction.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	rawTransaction.Value = valueBytes
	rawTransaction.ValueHash = valueBytesHash
	rawTransaction.ValueHashSignature = valueBytesHashSig

	return rawTransaction
}
