/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	tmos "github.com/ChengtayChain/ChengtayChain/libs/os"
	"github.com/ChengtayChain/ChengtayChain/libs/rand"
	"github.com/ChengtayChain/ChengtayChain/privval"
	rpcClient "github.com/ChengtayChain/ChengtayChain/rpc/client/http"
	ctypes "github.com/ChengtayChain/ChengtayChain/rpc/core/types"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// sendtxCmd represents the sendtx command
var SendtxCmd = &cobra.Command{
	Use:   "sendtx",
	Short: "Send a clearing Tx",
	Long:  `Send a clearing Tx`,
	RunE:  sendTx,
}
var privKey ed25519.PrivKeyEd25519
func initialize() {
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

}
func sendTx (cmd *cobra.Command, args []string) error{
	initialize()
	// get TMHOME

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

		time.Sleep(500 * time.Millisecond)
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
	value := types.IncomeTransactionValue{
		Timestamp: uint64(time.Now().Unix()),
		Items:     make([]types.IncomeTransactionItem, 0),
	}
	//nonce
	bytes256 := rand.Bytes(256)
	copy(value.Nonce[:256], bytes256[:256])

	//income 填入value，byte形式还是float
	vehicleOwnerIncomeInfo, totalPrice, chengtayPrice, _ := Compute("1")
	var IncomeTransactionItem = types.IncomeTransactionItem{
		*vehicleOwnerIncomeInfo,
		totalPrice,
		chengtayPrice,
	}
	value.Items = append(value.Items, IncomeTransactionItem)
	valueBytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	valueBytesHash := types.DefaultHashProvider.Digest(valueBytes)

	valueBytesHashSig, err := privKey.Sign(valueBytesHash)
	if err != nil {
		panic(err)
	}

	rawTransaction.Type = types.TransactionIncome
	rawTransaction.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	rawTransaction.Value = valueBytes
	rawTransaction.ValueHash = valueBytesHash
	rawTransaction.ValueHashSignature = valueBytesHashSig

	return rawTransaction
}
