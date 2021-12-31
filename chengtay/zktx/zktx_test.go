package zktx

import (
	"io/ioutil"
	"testing"
)

func TestVerifyProof(t *testing.T) {
	filepath := "/home/waris/circuit/proof"
	proof, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	VerifyProof(proof)
}
