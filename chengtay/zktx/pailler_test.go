package zktx

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"testing"
)

func TestMimcAll(t *testing.T) {
	assert := test.NewAssert(t)
	var paillerCircuit, witness PaillerCircuit
	// compile the circuit into a R1CS
	_, err := frontend.Compile(ecc.BN254, backend.GROTH16, &paillerCircuit)
	fmt.Println(err)
	assert.NoError(err)

	//// input
	//var N, Nsquare, R big.Int
	//var m, expectedResult big.Int
	//
	//m.SetString("10", 10)
	//expectedResult.SetString("1", 10)
	//N.SetString("3", 10)
	//Nsquare.SetString("9", 10)
	//R.SetString("2", 10)
	//var pk Publickey
	//
	//pk.NSquared.Assign(Nsquare)
	//pk.N.Assign(N)
	//witness.Data.Assign(m)
	//witness.Publickey = pk
	//witness.R.Assign(R)
	//witness.ExpectedResult.Assign(expectedResult)
	//assert.ProverSucceeded(&paillerCircuit, &witness, test.WithCurves(ecc.BN254))

	// input
	var pk = Publickey{
		N:        frontend.Value(3),
		NSquared: frontend.Value(9),
	}
	witness = PaillerCircuit{
		Publickey:      pk,
		R:              frontend.Value(2),
		Data:           frontend.Value(10),
		ExpectedResult: frontend.Value(4),
	}
	assert.ProverSucceeded(&paillerCircuit, &witness, test.WithCurves(ecc.BN254))
}
