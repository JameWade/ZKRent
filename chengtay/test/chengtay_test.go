package test

import (
	"encoding/json"
	"fmt"
	abciserver "github.com/ChengtayChain/ChengtayChain/abci/server"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	"github.com/ChengtayChain/ChengtayChain/libs/log"
	tmos "github.com/ChengtayChain/ChengtayChain/libs/os"
	"github.com/ChengtayChain/ChengtayChain/libs/rand"
	tmrand "github.com/ChengtayChain/ChengtayChain/libs/rand"
	"github.com/ChengtayChain/ChengtayChain/libs/service"
	"github.com/ChengtayChain/ChengtayChain/privval"
	rpcClient "github.com/ChengtayChain/ChengtayChain/rpc/client/http"
	"github.com/ChengtayChain/ChengtayChain/types/time"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	abcicli "github.com/ChengtayChain/ChengtayChain/abci/client"
	ctypes "github.com/ChengtayChain/ChengtayChain/rpc/core/types"
	//"sort"
	"testing"

	abcitypes "github.com/ChengtayChain/ChengtayChain/abci/types"
	chengtay "github.com/ChengtayChain/ChengtayChain/chengtay/abci"
)

const (
	testKey   = "abc"
	testValue = "def"
)
func RandVal(i int) abcitypes.ValidatorUpdate {
	pubkey := tmrand.Bytes(32)
	power := tmrand.Uint16() + 1
	v := abcitypes.Ed25519ValidatorUpdate(pubkey, int64(power))
	return v
}
func RandVals(cnt int) []abcitypes.ValidatorUpdate {
	res := make([]abcitypes.ValidatorUpdate, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = RandVal(i)
	}
	return res
}
//info() test
func TestChengtayInfo(t *testing.T) {
	dir, err := ioutil.TempDir("/home/runner", "abci-chengtay-test") // TODO
	if err != nil {
		t.Fatal(err)
	}
	chengtay := chengtay.NewApplication(dir)
	chengtay.InitChain(abcitypes.RequestInitChain{
		Validators: RandVals(1),
	})
	height := int64(0)

	resInfo := chengtay.Info(abcitypes.RequestInfo{})
	if resInfo.LastBlockHeight != height {
		t.Fatalf("expected height of %d, got %d", height, resInfo.LastBlockHeight)
	}

	// make and apply block
	height = int64(1)
	hash := []byte("foo")
	header := abcitypes.Header{
		Height: height,
	}
	chengtay.BeginBlock(abcitypes.RequestBeginBlock{Hash: hash, Header: header})
	chengtay.EndBlock(abcitypes.RequestEndBlock{Height: header.Height})
	chengtay.Commit()

	resInfo = chengtay.Info(abcitypes.RequestInfo{})
	if resInfo.LastBlockHeight != height {
		t.Fatalf("expected height of %d, got %d", height, resInfo.LastBlockHeight)
	}

}

// add a validator, remove a validator, update a validator
/* the func about validate doesn't complete
func TestValUpdates(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "abci-kvstore-test") // TODO
	if err != nil {
		t.Fatal(err)
	}
	chengtay := chengtay.NewApplication(dir)

	// init with some validators
	total := 10
	nInit := 5
	vals := RandVals(total)
	// iniitalize with the first nInit
	chengtay.InitChain(types.RequestInitChain{
		Validators: vals[:nInit],
	})

	vals1, vals2 := vals[:nInit], kvstore.Validators()
	valsEqual(t, vals1, vals2)

	var v1, v2, v3 types.ValidatorUpdate

	// add some validators
	v1, v2 = vals[nInit], vals[nInit+1]
	diff := []types.ValidatorUpdate{v1, v2}
	tx1 := MakeValSetChangeTx(v1.PubKey, v1.Power)
	tx2 := MakeValSetChangeTx(v2.PubKey, v2.Power)

	makeApplyBlock(t, kvstore, 1, diff, tx1, tx2)

	vals1, vals2 = vals[:nInit+2], kvstore.Validators()
	valsEqual(t, vals1, vals2)

	// remove some validators
	v1, v2, v3 = vals[nInit-2], vals[nInit-1], vals[nInit]
	v1.Power = 0
	v2.Power = 0
	v3.Power = 0
	diff = []types.ValidatorUpdate{v1, v2, v3}
	tx1 = MakeValSetChangeTx(v1.PubKey, v1.Power)
	tx2 = MakeValSetChangeTx(v2.PubKey, v2.Power)
	tx3 := MakeValSetChangeTx(v3.PubKey, v3.Power)

	makeApplyBlock(t, kvstore, 2, diff, tx1, tx2, tx3)

	vals1 = append(vals[:nInit-2], vals[nInit+1]) // nolint: gocritic
	vals2 = kvstore.Validators()
	valsEqual(t, vals1, vals2)

	// update some validators
	v1 = vals[0]
	if v1.Power == 5 {
		v1.Power = 6
	} else {
		v1.Power = 5
	}
	diff = []types.ValidatorUpdate{v1}
	tx1 = MakeValSetChangeTx(v1.PubKey, v1.Power)

	makeApplyBlock(t, kvstore, 3, diff, tx1)

	vals1 = append([]types.ValidatorUpdate{v1}, vals1[1:]...)
	vals2 = kvstore.Validators()
	valsEqual(t, vals1, vals2)

}
/*
func makeApplyBlock(
	t *testing.T,
	kvstore types.Application,
	heightInt int,
	diff []types.ValidatorUpdate,
	txs ...[]byte) {
	// make and apply block
	height := int64(heightInt)
	hash := []byte("foo")
	header := types.Header{
		Height: height,
	}

	kvstore.BeginBlock(types.RequestBeginBlock{Hash: hash, Header: header})
	for _, tx := range txs {
		if r := kvstore.DeliverTx(types.RequestDeliverTx{Tx: tx}); r.IsErr() {
			t.Fatal(r)
		}
	}
	resEndBlock := kvstore.EndBlock(types.RequestEndBlock{Height: header.Height})
	kvstore.Commit()

	valsEqual(t, diff, resEndBlock.ValidatorUpdates)

}

// order doesn't matter
func valsEqual(t *testing.T, vals1, vals2 []types.ValidatorUpdate) {
	if len(vals1) != len(vals2) {
		t.Fatalf("vals dont match in len. got %d, expected %d", len(vals2), len(vals1))
	}
	sort.Sort(types.ValidatorUpdates(vals1))
	sort.Sort(types.ValidatorUpdates(vals2))
	for i, v1 := range vals1 {
		v2 := vals2[i]
		if !bytes.Equal(v1.PubKey.Data, v2.PubKey.Data) ||
			v1.Power != v2.Power {
			t.Fatalf("vals dont match at index %d. got %X/%d , expect %X/%d", i, v2.PubKey, v2.Power, v1.PubKey, v1.Power)
		}
	}
}
*/




func makeGRPCClientServer(app abcitypes.Application) (abcicli.Client, service.Service, error) {
	// Start the listener
	var socket = "unix://127.0.0.1:26658"
	logger := log.TestingLogger()
	gapp := abcitypes.NewGRPCApplication(app)
	server := abciserver.NewGRPCServer(socket, gapp)
	server.SetLogger(logger.With("module", "abci-server"))
	if err := server.Start(); err != nil {
		return nil, nil, err
	}

	client := abcicli.NewGRPCClient(socket, true)
	client.SetLogger(logger.With("module", "abci-client"))
	if err := client.Start(); err != nil {
		server.Stop()
		return nil, nil, err
	}
	return client, server, nil
}

func TestClientServer(t *testing.T) {
	// set up grpc app
	chengTay := chengtay.NewApplication("/home/runner")
	gclient, gserver, err := makeGRPCClientServer(chengTay)
	require.Nil(t, err)
	defer gserver.Stop()
	defer gclient.Stop()

	testClient(t, gclient)
}




func testClient(t *testing.T,app abcicli.Client) {
	tx := GenTransaction()
	ar, err := app.DeliverTxSync(abcitypes.RequestDeliverTx{Tx: tx})
	require.NoError(t, err)
	require.Equal(t,3,int(ar.Code))
}


func TestClient(t *testing.T) {
	tx := GenTransaction()
	var address = "http://127.0.0.1:26657" // cfg.DefaultRPCConfig().ListenAddress
	var client *rpcClient.HTTP
	{
		var err error
		client, err = rpcClient.New(address, "/websocket")
		if err != nil {
			panic(err)
		}
	}

	ar, err := client.BroadcastTxSync(tx)
	fmt.Println(ar)
	require.NoError(t, err)
	require.Equal(t,0,int(ar.Code))
	// repeating tx doesn't raise error


	/* NO query
	// make sure query is fine
	//
	resQuery, err := app.QuerySync(types.RequestQuery{
		Path: "/store",
		Data: []byte(key),
	})
	require.Nil(t, err)
	require.Equal(t, code.CodeTypeOK, resQuery.Code)
	require.Equal(t, key, string(resQuery.Key))
	require.Equal(t, value, string(resQuery.Value))
	require.EqualValues(t, info.LastBlockHeight, resQuery.Height)

	// make sure proof is fine
	resQuery, err = app.QuerySync(types.RequestQuery{
		Path:  "/store",
		Data:  []byte(key),
		Prove: true,
	})
	require.Nil(t, err)
	require.Equal(t, code.CodeTypeOK, resQuery.Code)
	require.Equal(t, key, string(resQuery.Key))
	require.Equal(t, value, string(resQuery.Value))
	require.EqualValues(t, info.LastBlockHeight, resQuery.Height)

	 */
}

func init() {
	rand.Seed(time.Now().Unix())
}
func GenTransaction() []byte{
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

	rawTransaction := randomRawTransaction(privKey)
	bytes, err := json.Marshal(&rawTransaction)
	if err==nil{
		return bytes
	}else {
		return nil
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

func randomStorageItem() (item types.StorageItem) {
	item.CarID = types.ID(rand.Str(128))
	item.Timestamp = uint64(time.Now().Unix())
	item.ContentType = "whatever"
	item.Content = []byte(rand.Str(1 + rand.Intn(32768))) // https://github.com/tendermint/tendermint/pull/5215
	item.StorageItemID = types.ID(rand.Str(128))
	return item
}

func randomMerkleNode() (merkleNode types.IMerkleNode) {
	node := types.StorageItemMerkleNode(randomStorageItem())
	return &node
}

func randomMerkleTree() (merkleTree types.IMerkleTree) {
	merkleTree = &types.MerkleTree{}
	n := merkleTree.GetCapacity()
	for i := 0; i < n; i++ {
		err := merkleTree.SetMerkleNode(i, randomMerkleNode())
		if err != nil {
			panic(err)
		}
	}

	_, err := merkleTree.GetMerkleRoot()
	if err != nil {
		panic(err)
	}

	return merkleTree
}

func randomRawTransaction(privKey ed25519.PrivKeyEd25519) (rawTransaction types.RawTransaction) {
	treeNum := rand.Intn(20)
	trees := make([]types.IMerkleTree, 0)
	for i := 0; i < treeNum; i++ {
		trees = append(trees, randomMerkleTree())
	}

	value := types.MerkleRootTransactionValue{
		Timestamp: uint64(time.Now().Unix()),
		Items:     make([]types.MerkleRootTransactionItem, 0),
	}
	bytes256 := rand.Bytes(256)
	copy(value.Nonce[:256], bytes256[:256])

	for i := 0; i < treeNum; i++ {
		root, err := trees[i].GetMerkleRoot()
		if err != nil {
			panic(err)
		}
		value.Items = append(value.Items, root)
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

	rawTransaction.Type = types.TransactionMerkleroot
	rawTransaction.PublicKey = privKey.PubKey().(ed25519.PubKeyEd25519)
	rawTransaction.Value = valueBytes
	rawTransaction.ValueHash = valueBytesHash
	rawTransaction.ValueHashSignature = valueBytesHashSig

	return rawTransaction
}
