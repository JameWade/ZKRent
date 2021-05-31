package chengtay

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	abcitypes "github.com/ChengtayChain/ChengtayChain/abci/types"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/ChengtayChain/ChengtayChain/chengtay/zktx"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	"github.com/ChengtayChain/ChengtayChain/version"
	dbm "github.com/tendermint/tm-db"
	"math/big"
	"sort"
)

var ProtocolVersion version.Protocol = 0x1
var DatabaseVersion version.Protocol = 0x1
var requestTimeout = 3600
var percentage, _ = new(big.Int).SetString("0.2", 10)

type State struct {
	LastBlockHeight   int64 // The height of last block. Every time a new block is added, the height gets increased by 1.
	LastBlockHash     []byte
	TransactionCount  int64 // The number of transactions in all blocks.
	TrustedPublicKeys []ed25519.PubKeyEd25519

	RentAccount      map[string]*big.Int //pubkey: money
	LockRent         map[string]*big.Int //vehicle: money
	LockSig          map[string][]byte   //vehicle: signature = key
	VehicleUser      map[string]string   //  verhicle ID : user pubkey       //check if vehicle inused
	VehicleOwner     map[string]string   //vehicle ID : Owner publickey
	VehicleStartTime map[string]int64    //vehivle : time
	RequestTime      map[string]int64    //vehivle : time
	UntrustedUser    map[string]uint32   //user pubkey : untrust num
}

type TemporaryState struct {
	block map[string]*types.RawTransaction
	//copy the balance
	state     State
	blockTime int64
}

func newState() State {
	return State{
		LastBlockHash:     make([]byte, 0),
		LastBlockHeight:   0,
		TransactionCount:  0,
		TrustedPublicKeys: make([]ed25519.PubKeyEd25519, 0),
		RentAccount:       make(map[string]*big.Int),
		LockRent:          make(map[string]*big.Int),
		LockSig:           make(map[string][]byte),
		VehicleUser:       make(map[string]string),
		VehicleOwner:      make(map[string]string),
		VehicleStartTime:  make(map[string]int64),
		RequestTime:       make(map[string]int64), //vehivle : time
		UntrustedUser:     make(map[string]uint32),
	}
}

func newTemporaryState() TemporaryState {
	return TemporaryState{
		block: nil,
	}
}

func loadState(db dbm.DB) (state State, err error) {
	var loadExistingState bool
	{
		loadExistingState, err = db.Has([]byte("_DatabaseVersion"))
		if err != nil {
			return State{}, err
		}
	}

	if !loadExistingState {
		state = newState()
	} else {
		bytes, err := db.Get([]byte("_DatabaseVersion"))
		if err != nil {
			return State{}, err
		}

		if len(bytes) == 0 {
			panic(fmt.Errorf("The database may haven't been initialized."))
		}

		// check database version
		var ver uint64
		err = json.Unmarshal(bytes, &ver)
		if err != nil {
			return State{}, err
		}
		if ver != DatabaseVersion.Uint64() {
			return State{}, fmt.Errorf(fmt.Sprintf("Database version mismatch. Got %d, %d expected.", ver, DatabaseVersion.Uint64()))
		}

		// load the state
		{
			bytes, err := db.Get([]byte("_State"))
			if err != nil {
				return State{}, err
			}

			if len(bytes) == 0 {
				return State{}, fmt.Errorf(fmt.Sprintf("The database is corrupted. Key '%s' is expected.", "_State"))
			}

			// TODO: Performance: Replace JSON with Protobuf/Amino, to save disk space

			err = json.Unmarshal(bytes, &state)
			if err != nil {
				return State{}, fmt.Errorf(fmt.Sprintf("The database is corrupted. The value of key '%s' is supposed to be JSON, representing struct '%s'.%s", "_State", "State", err.Error()))
			}
		}

	}

	return state, nil
}

func saveState(state State, db dbm.DB) (err error) {
	// write database version
	{
		bytes, err := json.Marshal(DatabaseVersion)
		if err != nil {
			panic(err)
		}
		err = db.Set([]byte("_DatabaseVersion"), bytes)
		if err != nil {
			return err
		}
	}

	// write state data
	{
		bytes, err := json.Marshal(state)
		if err != nil {
			panic(err)
		}
		err = db.Set([]byte("_State"), bytes)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

var _ abcitypes.Application = (*Application)(nil)

type Application struct {
	abcitypes.BaseApplication

	db             dbm.DB
	state          State
	temporaryState TemporaryState
}

func NewApplication(dbDir string) *Application {
	name := "chengtay.data"

	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}

	state, err := loadState(db)
	if err != nil {
		panic(err)
	}

	return &Application{
		state:          state,
		temporaryState: newTemporaryState(),
		db:             db,
	}
}

//MerklerootTransaction need chengtay sign and must pass the whitelist
func (app *Application) parseAndCheckMerklerootTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	// parse
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}

	// check the signature
	// 1. the public key is in the whitelist
	// 2. the signature is valid
	fmt.Println("merkle VerifySignature start")
	// TODO: Performance: use a hash table/map to replace for-loop
	{
		transactionPublicKey := rawTransaction.PublicKey
		isTrusted := false
		for _, key := range app.state.TrustedPublicKeys {
			if key == transactionPublicKey {
				isTrusted = true
				break
			}
		}
		if !isTrusted {
			fmt.Println("ErrorPublicKeyUntrusted")
			return types.ErrorPublicKeyUntrusted, types.RawTransaction{}, nil
		}

		if !rawTransaction.VerifySignature() {
			return types.ErrorSignatureInvalid, types.RawTransaction{}, nil
		}

	}

	// all set
	return types.ErrorNoError, rawTransaction, nil
}

//only need verify signature do not need passwhitelist
func (app *Application) parseAndCheckTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	// parse
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	// check the signature
	//fmt.Println("transaction VerifySignature start")
	if !rawTransaction.VerifySignature() {
		fmt.Println("transaction VerifySignature not pass")
		return types.ErrorSignatureInvalid, types.RawTransaction{}, nil
	}
	// all set
	return types.ErrorNoError, rawTransaction, nil
}

func (app *Application) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Note: don't rely on CheckTx for anything
	var errorCode uint32
	var rawTransaction = types.RawTransaction{}
	err := json.Unmarshal(req.Tx, &rawTransaction)
	if err != nil {
		panic(err)
	}
	if rawTransaction.Type == types.TransactionMerkleroot {
		//MerklerootTransaction need chengtay sign and must pass the whitelist
		errorCode, _, err = app.parseAndCheckMerklerootTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else {
		//other transaction just verify the signature
		errorCode, _, err = app.parseAndCheckTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	}

	return abcitypes.ResponseCheckTx{Code: errorCode}
}

func (app *Application) parseAndDeliverRequestVehicleTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	//user should put a deposit on the chain,the deposit could use vehicle 1 month
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	RequestVehicleTransactionValue := types.RequestVehicleTransactionValue{}
	err = json.Unmarshal(rawTransaction.Value, &RequestVehicleTransactionValue)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}

	deposit := RequestVehicleTransactionValue.Deposit
	vehicleRequestPublicKey := RequestVehicleTransactionValue.PublicKey.String()
	vehicleId := RequestVehicleTransactionValue.VehicleID
	//if VehicleInused
	ifusedpublickey := app.temporaryState.state.VehicleUser[vehicleId]
	if  ifusedpublickey != ""{
		fmt.Println("**************************")
		fmt.Println("the vehicle Is being used ")
		fmt.Println("**************************")
		return types.ErrorVehicleInused, rawTransaction, nil
	}
	//if _, ok := app.state.LockSig[vehicleId]; ok {
	//	{
	//		fmt.Println("vehicle inused")
	//	}
	//	return types.ErrorVehicleInused, rawTransaction, nil
	//}
	fmt.Println(app.state.LockSig[vehicleId])
	if app.temporaryState.state.LockSig[vehicleId] != nil{
			{
				fmt.Println("**************************")
				fmt.Println("vehicle inused")
				fmt.Println("**************************")
			}
			return types.ErrorVehicleInused, rawTransaction, nil
	}

	//deal with the deposit,fixed deposit
	var monthRent *big.Int
	if app.temporaryState.state.UntrustedUser[vehicleRequestPublicKey] >  0{
		monthRent, _ = new(big.Int).SetString("200", 10)
	}else {
		monthRent, _ = new(big.Int).SetString("100", 10)
	}
	//fmt.Println(monthRent)

	if deposit.Cmp(monthRent) == -1 {
		//fmt.Println("balance not ok deliver")
		return types.ErrorInsufficientDeposit, types.RawTransaction{}, nil
	} else {
		//fmt.Println("balance ok deliver")
		//if no account on the chain,create an account
		if _, ok := app.temporaryState.state.RentAccount[vehicleRequestPublicKey]; !ok {
			//not exist this account,creat this account firstly
			app.temporaryState.state.RentAccount[vehicleRequestPublicKey], _ = new(big.Int).SetString("1000000", 10)
		}
		app.temporaryState.state.RentAccount[vehicleRequestPublicKey].Sub(app.temporaryState.state.RentAccount[vehicleRequestPublicKey], deposit)
		if app.temporaryState.state.RentAccount[vehicleRequestPublicKey].Cmp(deposit) != -1 {
			if _, ok := app.temporaryState.state.LockRent[vehicleId]; !ok {
				app.temporaryState.state.LockRent[vehicleId], _ = new(big.Int).SetString("0", 10)
			}
			app.temporaryState.state.LockRent[vehicleId].Add(app.temporaryState.state.LockRent[vehicleId], deposit)
			app.temporaryState.state.VehicleUser[vehicleId] = vehicleRequestPublicKey
			app.temporaryState.state.RequestTime[vehicleId] = app.temporaryState.blockTime
			fmt.Println("*****************************************************")
			fmt.Println("the vehicle request is being sent to the vehicleOwner")
			fmt.Println("*****************************************************")
		} else {
			fmt.Println("error")
			return types.ErrorInsufficientBalance, types.RawTransaction{}, nil
		}

	}
	return types.ErrorNoError, rawTransaction, nil
}

func (app *Application) parseAndDeliverResponseVehicleTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	ResponseVehicleTransactionValue := types.ResponseVehicleTransactionValue{}
	err = json.Unmarshal(rawTransaction.Value, &ResponseVehicleTransactionValue)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	//signature put into LockSig
	vehicleId := ResponseVehicleTransactionValue.VehicleID
	signature := ResponseVehicleTransactionValue.VehicleIDSignature
	publicKey := ResponseVehicleTransactionValue.PublicKey

	ifusedpublickey := app.temporaryState.state.VehicleUser[vehicleId]
	if  ifusedpublickey == ""{
		fmt.Println("****************************************")
		fmt.Println("the vehicle Is not being requested ")
		fmt.Println("****************************************")
		return types.ErrorVehicleInused, rawTransaction, nil
	}

	if app.state.LockSig[vehicleId] != nil{
		fmt.Println("the vehicle Is  being used ")
		return types.ErrorVehicleInused, rawTransaction, nil
	}
	app.temporaryState.state.LockSig[vehicleId] = signature
	app.temporaryState.state.VehicleOwner[vehicleId] = publicKey.String()
	fmt.Println("****************************************")
	fmt.Println(vehicleId)
	fmt.Println("the vehicle request processing completed")
	fmt.Println("****************************************")
	// start time
	app.temporaryState.state.VehicleStartTime[vehicleId] = app.temporaryState.blockTime
	return types.ErrorNoError, rawTransaction, nil
}

func (app *Application) parseAndDeliverResponseTimeoutTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	ResponseTimeoutTransactionValue := types.ResponseTimeoutTransactionValue{}
	err = json.Unmarshal(rawTransaction.Value, &ResponseTimeoutTransactionValue)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}

	vehicleId := ResponseTimeoutTransactionValue.VehicleID
	vehicleRequestPublicKey := ResponseTimeoutTransactionValue.PublicKey.String()

	//if signature exist ,then  Transaction invalid
	if _, ok := app.temporaryState.state.LockSig[vehicleId]; ok {
		fmt.Println("*****************************")
		fmt.Println("the vehicle Is  be signed ")
		fmt.Println("*****************************")
		return types.ErrorInvalidTimeout, rawTransaction, nil
	}
	ifusedpublickey := app.temporaryState.state.VehicleUser[vehicleId]
	if  ifusedpublickey != ""{
		fmt.Println("**************************")
		fmt.Println("the vehicle Is being used ")
		fmt.Println("**************************")
		return types.ErrorVehicleInused, rawTransaction, nil
	}

	if app.temporaryState.blockTime-app.temporaryState.state.RequestTime[vehicleId] > int64(requestTimeout) {
		app.temporaryState.state.RentAccount[vehicleRequestPublicKey].Add(app.temporaryState.state.RentAccount[vehicleRequestPublicKey], app.temporaryState.state.LockRent[vehicleId])
		app.temporaryState.state.LockRent[vehicleId].SetInt64(0)
		app.temporaryState.state.VehicleUser[vehicleId] = ""
		return types.ErrorNoError, rawTransaction, nil
	}

	return types.ErrorInvalidTimeout, rawTransaction, nil

}

func (app *Application) parseAndDeliverReturnVehicleTimeoutTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	//
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	ReturnVehicleTimeoutTransactionValue := types.ResponseVehicleTransactionValue{}
	err = json.Unmarshal(rawTransaction.Value, &ReturnVehicleTimeoutTransactionValue)
	if err != nil {
		return types.ErrorJsonParsing, types.RawTransaction{}, nil
	}
	vehicleId := ReturnVehicleTimeoutTransactionValue.VehicleID
	ifusedpublickey := app.temporaryState.state.VehicleUser[vehicleId]
	if  ifusedpublickey == ""{
		fmt.Println(vehicleId)
		fmt.Println("**********************************")
		fmt.Println("the vehicle was returned or unused")
		fmt.Println("**********************************")
		return types.ErrorVehicleUnused, rawTransaction, nil
	}

	//if signature exist ,then  Transaction invalid
	if app.temporaryState.state.LockSig[vehicleId] == nil {
		fmt.Println("**********************************")
		fmt.Println("the vehicle is not be signed")
		fmt.Println("**********************************")
		return types.ErrorVehicleUnused, rawTransaction, nil
	}

	latestTime := app.temporaryState.blockTime
	startTime := app.state.VehicleStartTime[vehicleId]
	useDays := int((latestTime - startTime) / (24 * 60 * 60))
	fmt.Println("latestTime:", latestTime)
	fmt.Println("startTime", startTime)
	fmt.Println("useDays", useDays)
	//todo : can not just use days,
	if useDays > 31 {
		UntrusedUserPublicKey := app.temporaryState.state.VehicleUser[vehicleId]
		app.temporaryState.state.UntrustedUser[UntrusedUserPublicKey] = app.temporaryState.state.UntrustedUser[UntrusedUserPublicKey] + 1   //count++
		//then
	}else {
		fmt.Println("**********************************")
		fmt.Println("used days < 31 days")
		fmt.Println("**********************************")
	}
	return types.ErrorNoError, rawTransaction, nil
}

func (app *Application) parseAndDeliverClearingTransaction(content []byte) (errorCode uint32, rawTransaction types.RawTransaction, criticalError error) {
	err := json.Unmarshal(content, &rawTransaction)
	if err != nil {
		return types.ErrorJsonParsing, rawTransaction, nil
	}
	clearingTransactionValue := types.ClearingTransactionValue{}
	err = json.Unmarshal(rawTransaction.Value, &clearingTransactionValue)
	if err != nil {
		return types.ErrorJsonParsing, rawTransaction, nil
	}

	vehicleId := clearingTransactionValue.VehicleID
	vehicleRequestPublicKey := app.temporaryState.state.VehicleUser[vehicleId]
	vehicleOwnerPublicKey := app.temporaryState.state.VehicleOwner[vehicleId]
	if  vehicleRequestPublicKey == ""{
		fmt.Println("**********************************************")
		fmt.Println(vehicleId)
		fmt.Println("Request is null ,the vehicle was not inused")
		fmt.Println("**********************************************")
		return types.ErrorVehicleUnused, rawTransaction, nil
	}
	//if signature exist ,then  Transaction invalid
	if app.temporaryState.state.LockSig[vehicleId] == nil {
		fmt.Println("the vehicle was been returned")
		return types.ErrorVehicleUnused, rawTransaction, nil
	}
	//usedays must > 30
	latestTime := app.temporaryState.blockTime
	startTime := app.state.VehicleStartTime[vehicleId]
	useDays := int((latestTime - startTime) / (24 * 60 * 60))
	fmt.Println("latestTime:", latestTime)
	fmt.Println("startTime", startTime)
	fmt.Println("useDays", useDays)
	//todo : can not just use days,
	if useDays < 30 {
		fmt.Println("**************************")
		fmt.Println("ErrorClearingDays")
		fmt.Println("**************************")
		return types.ErrorClearingDays,rawTransaction,nil
	}

	//compute   lockRent
	var monthRent *big.Int
	if _, ok := app.temporaryState.state.UntrustedUser[vehicleRequestPublicKey]; ok {
		monthRent, _ = new(big.Int).SetString("200", 10)
	} else {
		monthRent, _ = new(big.Int).SetString("100", 10)
	}

	ChengtayPublicKey := app.state.TrustedPublicKeys[0].String()
	chengtayPercentage := monthRent
	chengtayPercentage = chengtayPercentage.Mul(monthRent, percentage)

	app.temporaryState.state.RentAccount[ChengtayPublicKey] = app.temporaryState.state.RentAccount[ChengtayPublicKey].Add(chengtayPercentage, app.temporaryState.state.RentAccount[ChengtayPublicKey])
	ownerPercentage := monthRent
	ownerPercentage = ownerPercentage.Sub(monthRent, chengtayPercentage)
	app.temporaryState.state.RentAccount[vehicleOwnerPublicKey] = app.temporaryState.state.RentAccount[vehicleOwnerPublicKey].Add(chengtayPercentage, app.temporaryState.state.RentAccount[vehicleOwnerPublicKey])
	app.temporaryState.state.LockRent[vehicleId].Sub(app.temporaryState.state.LockRent[vehicleId], monthRent)
	//app.temporaryState.state.LockRent[vehicleId],_ = new(big.Int).SetString("0", 10)
	app.temporaryState.state.VehicleUser[vehicleId] = ""    //the vehicle is unused
	app.temporaryState.state.LockSig[vehicleId] = nil
	//delete(app.state.LockSig, vehicleId)
	fmt.Println("********************************")
	fmt.Println("the Vehicle has been returned")
	fmt.Println("********************************")
	fmt.Println(app.temporaryState.state.RentAccount[ChengtayPublicKey])


	return types.ErrorNoError, rawTransaction, nil
}

func (app *Application) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	// Note: don't rely on CheckTx for anything,so rechecktx
	var errorCode uint32
	var rawTransaction = types.RawTransaction{}
	err := json.Unmarshal(req.Tx, &rawTransaction)
	if err != nil {
		panic(err)
	}
	if rawTransaction.Type == types.TransactionMerkleroot {
		//MerklerootTransaction need chengtay sign and must pass the whitelist
		errorCode, _, err = app.parseAndCheckMerklerootTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else {
		//other transaction just verify the signature
		errorCode, _, err = app.parseAndCheckTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	}
	if errorCode != types.ErrorNoError {
		return abcitypes.ResponseDeliverTx{Code: errorCode}
	}

	//deliver TransactionRequestVehicle and others
	if rawTransaction.Type == types.TransactionRequestVehicle {
		//TransactionRequestVehicle need a deposit and deal with the account
		errorCode, _, err = app.parseAndDeliverRequestVehicleTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else if rawTransaction.Type == types.TransactionResponseVehicle {
		errorCode, _, err = app.parseAndDeliverResponseVehicleTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else if rawTransaction.Type == types.TransactionResponseTimeout {
		errorCode, _, err = app.parseAndDeliverResponseTimeoutTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else if rawTransaction.Type == types.TransactionReturnVehicleTimeout {
		errorCode, _, err = app.parseAndDeliverReturnVehicleTimeoutTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	} else if rawTransaction.Type == types.TransactionClearing {
		errorCode, _, err = app.parseAndDeliverClearingTransaction(req.Tx)
		if err != nil {
			panic(err)
		}
	}

	if errorCode != types.ErrorNoError {
		return abcitypes.ResponseDeliverTx{Code: errorCode}
	}

	//write to the chain
	errorCode, transaction, err := types.GetTransaction(rawTransaction)
	if err != nil {
		panic(err)
	}
	if errorCode != types.ErrorNoError {
		return abcitypes.ResponseDeliverTx{Code: errorCode}
	}
	var transactionHash string
	switch transaction.Type {
	case types.TransactionMerkleroot, types.TransactionRequestVehicle, types.TransactionResponseVehicle,
		types.TransactionResponseTimeout, types.TransactionReturnVehicleTimeout, types.TransactionClearing:
		hash, err := transaction.Value.GetHash()
		if err != nil {
			panic(err)
		}
		transactionHash = string(hash)
	default:
		panic(fmt.Errorf("Assert failed. Unsupported transaction type %d, where app.parseAndCheckTransaction() is supposed to report error code %d.", transaction.Type, types.ErrorUnknownTransactionType))
	}

	app.temporaryState.block[transactionHash] = &rawTransaction

	return abcitypes.ResponseDeliverTx{Code: types.ErrorNoError}
}

func (app *Application) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	if app.temporaryState.block != nil {
		panic(fmt.Errorf("Assert failed: The last block hasn't been committed yet."))
	}

	app.temporaryState.block = make(map[string]*types.RawTransaction)
	//get the time
	app.temporaryState.blockTime = req.Header.Time.Unix()
	//copy the balance
	//TODO: just
	bytes, err := json.Marshal(app.state)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &app.temporaryState.state)
	if err != nil {
		panic(err)
	}
	return abcitypes.ResponseBeginBlock{}
}

func (app *Application) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func BlockHash(lastBlockHeight int64, lastBlockHash []byte, transactionHashes []string) (ret []byte) {
	// sort the transactionHashes
	sort.Strings(transactionHashes)

	bytes := make([]byte, 0)

	// 1. the json byte of lastBlockHeight
	{
		jsonBytes, err := json.Marshal(lastBlockHeight)
		if err != nil {
			panic(err)
		}
		bytes = append(bytes, jsonBytes...)
	}

	// 2. lastBlockHash
	bytes = append(bytes, lastBlockHash...)

	// 3. transactionHashes (sorted)
	for _, transactionHash := range transactionHashes {
		bytes = append(bytes, []byte(transactionHash)...)
	}

	return types.DefaultHashProvider.Digest(bytes)
}

func (app *Application) Commit() abcitypes.ResponseCommit {
	bytes, err := json.Marshal(app.temporaryState.state)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &app.state)
	if err != nil {
		panic(err)
	}

	if app.temporaryState.block == nil {
		panic(fmt.Errorf("Assert failed: The current block hasn't been beginned yet."))
	}

	// save transactions to database
	for transactionHash, transaction := range app.temporaryState.block {
		key := []byte("transaction-" + hex.EncodeToString([]byte(transactionHash)))
		bytes, err := app.db.Get(key)
		if err != nil {
			panic(err)
		}

		if len(bytes) != 0 {
			panic(fmt.Errorf(fmt.Sprintf("Database corrupted. Transaction (Hash=%s) already exists in the database.", hex.EncodeToString([]byte(transactionHash)))))
		}

		bytes, err = json.Marshal(&transaction)
		err = app.db.Set(key, bytes)
		if err != nil {
			panic(err)
		}
	}
	// calculate block hash
	transactionHashes := make([]string, 0, len(app.temporaryState.block))
	for transactionHash := range app.temporaryState.block {
		transactionHashes = append(transactionHashes, transactionHash)
	}
	blockHash := BlockHash(app.state.LastBlockHeight, app.state.LastBlockHash, transactionHashes)
	// update state
	app.state.LastBlockHeight++
	app.state.LastBlockHash = blockHash
	app.state.TransactionCount += int64(len(app.temporaryState.block))

	app.temporaryState.block = nil

	// save state to the database
	err = saveState(app.state, app.db)
	if err != nil {
		panic(err)
	}
	//fmt.Println(app.state)
	return abcitypes.ResponseCommit{Data: app.state.LastBlockHash}
}

func (app *Application) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	topValidator := req.Validators[0]
	if topValidator.PubKey.Type == abcitypes.PubKeyEd25519 {
		var bytes32 [ed25519.PubKeyEd25519Size]byte
		// assert len(topValidator.PubKey.Data) == ed25519.PubKeyEd25519Size
		for k, v := range topValidator.PubKey.Data {
			bytes32[k] = v
		}

		//fmt.Println("DEBUG: Reading public key:", hex.EncodeToString(topValidator.PubKey.Data))

		publicKey := ed25519.PubKeyEd25519(bytes32)

		//fmt.Println("DEBUG: Trusted address:", publicKey.Address())

		app.state.TrustedPublicKeys = append(app.state.TrustedPublicKeys, publicKey)
		app.state.RentAccount[publicKey.String()], _ = new(big.Int).SetString("0", 10)
		//fmt.Println("chengtay is ",ChengtayPublicKey)
		//fmt.Println("\n\n\n\n\n")
	} else {
		panic(fmt.Errorf("Unsupported public key type."))
	}

	return abcitypes.ResponseInitChain{}
}

func (app *Application) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	data := make(map[string]interface{})
	{
		var err error
		data["State"], err = json.Marshal(app.state)
		if err != nil {
			panic(err)
		}
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return abcitypes.ResponseInfo{
		Data:             string(dataBytes),
		Version:          version.ABCIVersion,
		AppVersion:       uint64(ProtocolVersion),
		LastBlockHeight:  app.state.LastBlockHeight,
		LastBlockAppHash: app.state.LastBlockHash,
	}
}

//func (app *Application)query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
//	if reqQuery.Prove {
//		resQuery.Key = reqQuery.Data
//		pubkey := ed25519.PubKeyEd25519{}
//		copy(pubkey[:32],reqQuery.Data)
//		resQuery.Value = []byte(app.temporaryState.state.LockRent[pubkey.String()].String())
//		resQuery.Height = app.state.LastBlockHeight
//
//		return
//	}
//	resQuery.Key = reqQuery.Data
//	pubkey := ed25519.PubKeyEd25519{}
//	copy(pubkey[:32],reqQuery.Data)
//	resQuery.Value = []byte(app.temporaryState.state.LockRent[pubkey.String()].String())
//	resQuery.Height = app.state.LastBlockHeight
//	return resQuery
//}

func VerifyProof(tx *types.Transaction) bool {
	if err := zktx.VerifySendProof(tx.ZKSN(), tx.ZKCMTS(), tx.ZKProof(), &cmtbalance, tx.ZKCMT()); err != nil {
		fmt.Println("invalid zk send proof: ", err)
		return false
	}
	return true
}
