package zktx

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"testing"
)

func TestModAll(t *testing.T) {
	assert := test.NewAssert(t)
	var modCircuit, witness ModCircuit
	// compile the circuit into a R1CS
	_, err := frontend.Compile(ecc.BN254, backend.GROTH16, &modCircuit)
	fmt.Println(err)
	assert.NoError(err)
	witness = ModCircuit{
		X:              frontend.Value(30),
		Y:              frontend.Value(3),
		ExpectedResult: frontend.Value(30),
	}
	assert.ProverSucceeded(&modCircuit, &witness, test.WithCurves(ecc.BN254))

	//big.Int{}
}
