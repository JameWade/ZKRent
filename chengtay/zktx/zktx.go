package zktx

// #cgo CFLAGS:  -I./include
// #cgo LDFLAGS: -L/home/waris/ZKRent/chengtay/zktx -L/usr/local/lib -lzkRent_verify -lff  -lsnark -lstdc++  -lgmp -lgmpxx -lzk_CircuitReader -lzk_Util
//#include"zkRent_verify.hpp"
//#include <stdlib.h>
import "C"
import "unsafe"

func VerifyProof(proof []byte) bool {

	tf := C.verifyproof((*C.char)(unsafe.Pointer(&proof[0])))
	if tf == 1 {
		return true
	}
	return false

}
