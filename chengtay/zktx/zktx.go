package zktx

/*
#cgo LDFLAGS: -L/usr/local/lib -lzk_mint  -lzk_send  -lzk_deposit -lzk_redeem -lff  -lsnark -lstdc++  -lgmp -lgmpxx -lzkRent_verify
#include "mintcgo.hpp"
#include "sendcgo.hpp"
#include "depositcgo.hpp"
#include "redeemcgo.hpp"
#include "zkRent_verify.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"io"
	"math/big"
	"os"
	"sync"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Sequence struct {
	SN     *common.Hash
	CMT    *common.Hash
	Random *common.Hash
	Value  uint64
	Valid  bool
	Lock   sync.Mutex
}

type WriteSn struct {
	SNumber      *Sequence
	SNumberAfter *Sequence
}
type SequenceS struct {
	Suquence1 Sequence
	Suquence2 Sequence
	SNS       *Sequence
	PKBX      *big.Int
	PKBY      *big.Int
	Stage     uint8
}

const (
	Origin = iota
	Mint
	Send
	Update
	Deposit
	Redeem
)

var SNfile *os.File
var FileLine uint8

var Stage uint8
var ZKTxAddress = common.HexToAddress("ffffffffffffffffffffffffffffffffffffffff")
var SequenceNumber = InitializeSN()                //--zy
var SequenceNumberAfter *Sequence = InitializeSN() //--zy
var SNS *Sequence = nil

var ZKCMTNODES = 1 // max is 32  because of merkle leaves in libnsark is 32

var ErrSequence = errors.New("invalid sequence")
var RandomReceiverPK *ecdsa.PublicKey = nil

func InitializeSN() *Sequence {
	// For large-scale test, we suppose that SK = CRH(addr), there is impossible in pratical.
	SK := ZKTxAddress.Hash()
	r := &common.Hash{}
	sn := ComputePRF(SK.Bytes(), r.Bytes()) // sn = PRF(sk, r)

	cmt := GenCMT(0, sn.Bytes(), r.Bytes())
	return &Sequence{
		SN:     sn,
		CMT:    cmt,
		Random: r,
		Value:  0,
	}
}

func NewRandomHash() *common.Hash {
	uuid := make([]byte, 32)
	io.ReadFull(rand.Reader, uuid)
	hash := common.BytesToHash(uuid)
	return &hash
}

func NewRandomAddress() *common.Address {
	uuid := make([]byte, 20)
	io.ReadFull(rand.Reader, uuid)
	addr := common.BytesToAddress(uuid)
	return &addr
}

func NewRandomInt() *big.Int {
	uuid := make([]byte, 32)
	io.ReadFull(rand.Reader, uuid)
	r := big.NewInt(0)
	r.SetBytes(uuid)
	return r
}



var InvalidMintProof = errors.New("Verifying mint proof failed!!!")

func VerifyMintProof(cmtold *common.Hash, snaold *common.Hash, cmtnew *common.Hash, value uint64, proof []byte) error {
	cproof := C.CString(string(proof))
	cmtA_old_c := C.CString(hexutil.Encode(cmtold[:]))
	cmtA_c := C.CString(hexutil.Encode(cmtnew[:]))
	sn_old_c := C.CString(hexutil.Encode(snaold.Bytes()[:]))
	value_s_c := C.ulong(value)
	tf := C.verifyMintproof(cproof, cmtA_old_c, sn_old_c, cmtA_c, value_s_c)
	if tf == false {
		return InvalidMintProof
	}
	return nil
}

var InvalidSendProof = errors.New("Verifying send proof failed!!!")


func VerifySendProof(sna *common.Hash, cmts *common.Hash, proof []byte, cmtAold *common.Hash, cmtAnew *common.Hash) error {
	cproof := C.CString(string(proof))
	snAold_c := C.CString(hexutil.Encode(sna.Bytes()[:]))
	cmtS := C.CString(hexutil.Encode(cmts[:]))
	cmtAold_c := C.CString(hexutil.Encode(cmtAold[:]))
	cmtAnew_c := C.CString(hexutil.Encode(cmtAnew[:]))

	tf := C.verifySendproof(cproof, cmtAold_c, snAold_c, cmtS, cmtAnew_c)
	if tf == false {
		return InvalidSendProof
	}
	return nil
}
func Verifyproof() error {
	tf := C.verifyproof()
	if tf == false {
		return InvalidSendProof
	}
	return nil
}

var InvalidUpdateProof = errors.New("Verifying update proof failed!!!")



var InvalidDepositProof = errors.New("Verifying Deposit proof failed!!!")

func VerifyDepositProof(pk_recv *ecdsa.PublicKey, rtcmt common.Hash, cmtb *common.Hash, snb *common.Hash, cmtbnew *common.Hash, sns *common.Hash, proof []byte) error {
	PK_recv := crypto.PubkeyToAddress(*pk_recv) //--zy
	pk_recv_c := C.CString(hexutil.Encode(PK_recv[:]))
	cproof := C.CString(string(proof))
	rtmCmt := C.CString(hexutil.Encode(rtcmt[:]))
	cmtB := C.CString(hexutil.Encode(cmtb[:]))
	cmtBnew := C.CString(hexutil.Encode(cmtbnew[:]))
	SNB_c := C.CString(hexutil.Encode(snb.Bytes()[:]))
	SNS_c := C.CString(hexutil.Encode(sns.Bytes()[:]))
	tf := C.verifyDepositproof(cproof, rtmCmt, pk_recv_c, cmtB, SNB_c, cmtBnew, SNS_c)
	if tf == false {
		return InvalidDepositProof
	}
	return nil
}

var InvalidRedeemProof = errors.New("Verifying redeem proof failed!!!")

func VerifyRedeemProof(cmtold *common.Hash, snaold *common.Hash, cmtnew *common.Hash, value uint64, proof []byte) error {
	cproof := C.CString(string(proof))
	cmtA_old_c := C.CString(hexutil.Encode(cmtold[:]))
	cmtA_c := C.CString(hexutil.Encode(cmtnew[:]))
	sn_old_c := C.CString(hexutil.Encode(snaold.Bytes()[:]))
	value_s_c := C.ulong(value)

	tf := C.verifyRedeemproof(cproof, cmtA_old_c, sn_old_c, cmtA_c, value_s_c)
	if tf == false {
		return InvalidRedeemProof
	}
	return nil
}



//GenCMT生成CMT 调用c的sha256函数  （go的sha256函数与c有一些区别）
func GenCMT(value uint64, sn []byte, r []byte) *common.Hash {
	//sn_old_c := C.CString(hexutil.Encode(SNold[:]))
	value_c := C.ulong(value)
	sn_string := hexutil.Encode(sn[:])
	sn_c := C.CString(sn_string)
	defer C.free(unsafe.Pointer(sn_c))
	r_string := hexutil.Encode(r[:])
	r_c := C.CString(r_string)
	defer C.free(unsafe.Pointer(r_c))

	cmtA_c := C.genCMT(value_c, sn_c, r_c)
	cmtA_go := C.GoString(cmtA_c)
	//res := []byte(cmtA_go)
	res, _ := hex.DecodeString(cmtA_go)
	reshash := common.BytesToHash(res)
	return &reshash
}

//GenCMT生成CMT 调用c的sha256函数  （go的sha256函数与c有一些区别）
func GenCMTS(values uint64, pk_recv *ecdsa.PublicKey, rs []byte, sna []byte) *common.Hash {
	values_c := C.ulong(values)
	PK_recv := crypto.PubkeyToAddress(*pk_recv) //--zy
	pk_recv_c := C.CString(hexutil.Encode(PK_recv[:]))
	rs_string := hexutil.Encode(rs[:])
	rs_c := C.CString(rs_string)
	defer C.free(unsafe.Pointer(rs_c))
	sna_string := hexutil.Encode(sna[:])
	sna_c := C.CString(sna_string)
	defer C.free(unsafe.Pointer(sna_c))
	//uint64_t value_s,char* pk_string,char* sn_s_string,char* r_s_string,char *sn_old_string
	cmtA_c := C.genCMTS(values_c, pk_recv_c, rs_c, sna_c) //64长度16进制数
	cmtA_go := C.GoString(cmtA_c)
	//res := []byte(cmtA_go)
	res, _ := hex.DecodeString(cmtA_go)
	reshash := common.BytesToHash(res) //32长度byte数组
	return &reshash
}

//ComputePRF生成sn 调用c的sha256函数  （go的sha256函数与c有一些区别）
func ComputePRF(sk []byte, r []byte) *common.Hash {
	addr_string := hexutil.Encode(sk[:])
	addr_c := C.CString(addr_string)
	defer C.free(unsafe.Pointer(addr_c))

	r_string := hexutil.Encode(r[:])
	r_c := C.CString(r_string)
	defer C.free(unsafe.Pointer(r_c))

	sn_c := C.computePRF(addr_c, r_c)
	sn_go := C.GoString(sn_c)
	//res := []byte(cmtA_go)
	res, _ := hex.DecodeString(sn_go)
	reshash := common.BytesToHash(res)
	return &reshash
}

//ComputeCRH生成r_s 调用c的sha256函数  （go的sha256函数与c有一些区别）
func ComputeCRH(pk_recv common.Address, r []byte) *common.Hash {
	//PK_recv := crypto.PubkeyToAddress(*pk_recv) //--zy
	pk_recv_c := C.CString(hexutil.Encode(pk_recv[:]))

	r_string := hexutil.Encode(r[:])
	r_c := C.CString(r_string)
	defer C.free(unsafe.Pointer(r_c))

	r_s_c := C.computeCRH(pk_recv_c, r_c)
	r_s_go := C.GoString(r_s_c)
	//res := []byte(cmtA_go)
	res, _ := hex.DecodeString(r_s_go)
	reshash := common.BytesToHash(res)
	return &reshash
}

//GenRT 返回merkel树的hash  --zy
func GenRT(CMTSForMerkle []*common.Hash) common.Hash {
	var cmtArray string
	for i := 0; i < len(CMTSForMerkle); i++ {
		s := string(hexutil.Encode(CMTSForMerkle[i][:]))
		cmtArray += s
	}
	cmtsM := C.CString(cmtArray)
	rtC := C.genRoot(cmtsM, C.int(len(CMTSForMerkle))) //--zy
	rtGo := C.GoString(rtC)

	res, _ := hex.DecodeString(rtGo)   //返回32长度 []byte  一个byte代表两位16进制数
	reshash := common.BytesToHash(res) //32长度byte数组
	return reshash
}

func ComputeR(sk *big.Int) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{} //tbd
}



type AUX struct {
	Value uint64
	//SNs   *common.Hash
	Rs    *common.Hash
	SNa   *common.Hash
}


func GenerateKeyForRandomB(R *ecdsa.PublicKey, kB *ecdsa.PrivateKey) *ecdsa.PrivateKey {
	//skB*R
	c := kB.PublicKey.Curve
	tx, ty := c.ScalarMult(R.X, R.Y, kB.D.Bytes())
	tmp := tx.Bytes()
	tmp = append(tmp, ty.Bytes()...)

	//生成hash值H(skB*R)
	h := sha256.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	bs[0] = bs[0] % 128
	i := new(big.Int)
	i = i.SetBytes(bs)

	//生成公钥
	sx, sy := c.ScalarBaseMult(bs)
	sskB := new(ecdsa.PrivateKey)
	sskB.PublicKey.X, sskB.PublicKey.Y = c.Add(sx, sy, kB.PublicKey.X, kB.PublicKey.Y)
	sskB.Curve = c
	//生成私钥
	sskB.D = i.Add(i, kB.D)
	return sskB
}

func GenMintProof(ValueOld uint64, RAold *common.Hash, SNAnew *common.Hash, RAnew *common.Hash, CMTold *common.Hash, SNold *common.Hash, CMTnew *common.Hash, ValueNew uint64, SK *common.Hash) []byte {
	value_c := C.ulong(ValueNew)     //转换后零知识余额对应的明文余额
	value_old_c := C.ulong(ValueOld) //转换前零知识余额对应的明文余额

	sn_old_c := C.CString(hexutil.Encode(SNold[:]))
	r_old_c := C.CString(hexutil.Encode(RAold[:]))
	sn_c := C.CString(hexutil.Encode(SNAnew[:]))
	r_c := C.CString(hexutil.Encode(RAnew[:]))

	cmtA_old_c := C.CString(hexutil.Encode(CMTold[:])) //对于CMT  需要将每一个byte拆为两个16进制字符
	cmtA_c := C.CString(hexutil.Encode(CMTnew[:]))

	value_s_c := C.ulong(ValueNew - ValueOld) //需要被转换的明文余额

	sk_c := C.CString(hexutil.Encode(SK[:]))

	cproof := C.genMintproof(value_c, value_old_c, sn_old_c, r_old_c, sn_c, r_c, cmtA_old_c, cmtA_c, value_s_c, sk_c)

	var goproof string
	goproof = C.GoString(cproof)
	return []byte(goproof)
}

func GenSendProof(CMTA *common.Hash, ValueA uint64, RA *common.Hash, ValueS uint64, pk_recv *ecdsa.PublicKey, RS *common.Hash, SNA *common.Hash, CMTS *common.Hash, ValueAnew uint64, SNAnew *common.Hash, RAnew *common.Hash, CMTAnew *common.Hash, SK *common.Hash, pk_sender common.Address) []byte {
	cmtA_c := C.CString(hexutil.Encode(CMTA[:]))
	valueA_c := C.ulong(ValueA)
	rA_c := C.CString(hexutil.Encode(RA.Bytes()[:]))
	valueS := C.ulong(ValueS)
	PK_recv := crypto.PubkeyToAddress(*pk_recv) //--zy
	pk_recv_c := C.CString(hexutil.Encode(PK_recv[:]))
	rS := C.CString(hexutil.Encode(RS.Bytes()[:]))
	snA := C.CString(hexutil.Encode(SNA.Bytes()[:]))
	cmtS := C.CString(hexutil.Encode(CMTS[:]))
	//ValueAnew uint64 , SNAnew *common.Hash, RAnew *common.Hash,CMTAnew *common.Hash
	valueANew_c := C.ulong(ValueAnew)
	snAnew_c := C.CString(hexutil.Encode(SNAnew.Bytes()[:]))
	rAnew_c := C.CString(hexutil.Encode(RAnew.Bytes()[:]))
	cmtAnew_c := C.CString(hexutil.Encode(CMTAnew[:]))

	sk_c := C.CString(hexutil.Encode(SK[:]))
	//PK_sender := crypto.PubkeyToAddress(*pk_sender) //--zy
	pk_sender_c := C.CString(hexutil.Encode(pk_sender[:]))

	cproof := C.genSendproof(valueA_c, rS, snA, rA_c, cmtS, cmtA_c, valueS, pk_recv_c, valueANew_c, snAnew_c, rAnew_c, cmtAnew_c, sk_c, pk_sender_c)
	var goproof string
	goproof = C.GoString(cproof)
	return []byte(goproof)
}


func GenDepositProof(CMTS *common.Hash, ValueS uint64, SNS *common.Hash, RS *common.Hash, SNA *common.Hash, ValueB uint64, RB *common.Hash, SNBnew *common.Hash, RBnew *common.Hash, pk_recv *ecdsa.PublicKey, RTcmt []byte, CMTB *common.Hash, SNB *common.Hash, CMTBnew *common.Hash, CMTSForMerkle []*common.Hash, SK *common.Hash) []byte {
	cmtS_c := C.CString(hexutil.Encode(CMTS[:]))
	valueS_c := C.ulong(ValueS)
	PK_recv := crypto.PubkeyToAddress(*pk_recv) //--zy
	pk_recv_c := C.CString(hexutil.Encode(PK_recv[:]))
	SNS_c := C.CString(hexutil.Encode(SNS.Bytes()[:])) //--zy
	RS_c := C.CString(hexutil.Encode(RS.Bytes()[:]))   //--zy
	SNA_c := C.CString(hexutil.Encode(SNA.Bytes()[:]))
	valueB_c := C.ulong(ValueB)
	RB_c := C.CString(hexutil.Encode(RB.Bytes()[:])) //rA_c := C.CString(string(RA.Bytes()[:]))
	SNB_c := C.CString(hexutil.Encode(SNB.Bytes()[:]))
	SNBnew_c := C.CString(hexutil.Encode(SNBnew.Bytes()[:]))
	RBnew_c := C.CString(hexutil.Encode(RBnew.Bytes()[:]))
	cmtB_c := C.CString(hexutil.Encode(CMTB[:]))
	RT_c := C.CString(hexutil.Encode(RTcmt)) //--zy   rt

	cmtBnew_c := C.CString(hexutil.Encode(CMTBnew[:]))
	valueBNew_c := C.ulong(ValueB + ValueS)

	SK_c := C.CString(hexutil.Encode(SK.Bytes()[:]))

	var cmtArray string
	for i := 0; i < len(CMTSForMerkle); i++ {
		s := string(hexutil.Encode(CMTSForMerkle[i][:]))
		cmtArray += s
	}
	cmtsM := C.CString(cmtArray)
	nC := C.int(len(CMTSForMerkle))

	cproof := C.genDepositproof(valueBNew_c, valueB_c, SNB_c, RB_c, SNBnew_c, RBnew_c, SNS_c, RS_c, cmtB_c, cmtBnew_c, valueS_c, pk_recv_c, SNA_c, cmtS_c, cmtsM, nC, RT_c, SK_c)
	var goproof string
	goproof = C.GoString(cproof)
	return []byte(goproof)
}

func GenRedeemProof(ValueOld uint64, RAold *common.Hash, SNAnew *common.Hash, RAnew *common.Hash, CMTold *common.Hash, SNold *common.Hash, CMTnew *common.Hash, ValueNew uint64, SK *common.Hash) []byte {
	value_c := C.ulong(ValueNew)     //转换后零知识余额对应的明文余额
	value_old_c := C.ulong(ValueOld) //转换前零知识余额对应的明文余额

	sn_old_c := C.CString(hexutil.Encode(SNold.Bytes()[:]))
	r_old_c := C.CString(hexutil.Encode(RAold.Bytes()[:]))
	sn_c := C.CString(hexutil.Encode(SNAnew.Bytes()[:]))
	r_c := C.CString(hexutil.Encode(RAnew.Bytes()[:]))

	cmtA_old_c := C.CString(hexutil.Encode(CMTold[:])) //对于CMT  需要将每一个byte拆为两个16进制字符
	cmtA_c := C.CString(hexutil.Encode(CMTnew[:]))

	value_s_c := C.ulong(ValueOld - ValueNew) //需要被转换的明文余额

	SK_c := C.CString(hexutil.Encode(SK.Bytes()[:]))

	cproof := C.genRedeemproof(value_c, value_old_c, sn_old_c, r_old_c, sn_c, r_c, cmtA_old_c, cmtA_c, value_s_c, SK_c)

	var goproof string
	goproof = C.GoString(cproof)
	return []byte(goproof)
}

func GenR() *ecdsa.PrivateKey {
	Ka, err := crypto.GenerateKey()
	if err != nil {
		return nil
	}
	return Ka
}

func NewRandomPubKey(sA *big.Int, pkB ecdsa.PublicKey) *ecdsa.PublicKey {
	//sA*pkB
	c := pkB.Curve
	tx, ty := c.ScalarMult(pkB.X, pkB.Y, sA.Bytes())
	tmp := tx.Bytes()
	tmp = append(tmp, ty.Bytes()...)

	//生成hash值H(sA*pkB)
	h := sha256.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	bs[0] = bs[0] % 128

	//生成用于加密的公钥H(sA*pkB)P+pkB
	sx, sy := c.ScalarBaseMult(bs)
	spkB := new(ecdsa.PublicKey)
	spkB.X, spkB.Y = c.Add(sx, sy, pkB.X, pkB.Y)
	spkB.Curve = c
	return spkB
}



///zkRent
func VerifyReqProof(sna *common.Hash, cmts *common.Hash, proof []byte, cmtAold *common.Hash, cmtAnew *common.Hash) error {
	cproof := C.CString(string(proof))
	snAold_c := C.CString(hexutil.Encode(sna.Bytes()[:]))
	cmtS := C.CString(hexutil.Encode(cmts[:]))
	cmtAold_c := C.CString(hexutil.Encode(cmtAold[:]))
	cmtAnew_c := C.CString(hexutil.Encode(cmtAnew[:]))

	tf := C.verifySendproof(cproof, cmtAold_c, snAold_c, cmtS, cmtAnew_c)
	if tf == false {
		return InvalidSendProof
	}
	return nil
}
