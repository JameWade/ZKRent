package zktx

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

type PaillerCircuit struct {
	Publickey      Publickey
	R              frontend.Variable
	Data           frontend.Variable
	ExpectedResult frontend.Variable `gnark:"data,public"`
}
type Publickey struct {
	N frontend.Variable
	//G        frontend.Variable
	NSquared frontend.Variable
}

func (paillerCircuit *PaillerCircuit) mod(v1 frontend.Variable, v2 frontend.Variable, cs frontend.API) frontend.Variable {
	dis := cs.Div(v1, v2)
	res := cs.Sub(v1, cs.Mul(v2, dis))
	fmt.Println(res.WitnessValue)
	//mod = v1-dis*v2
	return res
}

func (paillerCircuit *PaillerCircuit) exp(cs frontend.API) frontend.Variable {
	// number of bits of exponent
	const bitSize = 8

	// specify constraints
	output := cs.Constant(1)
	bits := cs.ToBinary(paillerCircuit.Publickey.N, bitSize)
	cs.ToBinary(paillerCircuit.Publickey.N, bitSize)

	for i := 0; i < len(bits); i++ {

		if i != 0 {
			output = cs.Mul(output, output)
		}
		multiply := cs.Mul(output, paillerCircuit.R)
		output = cs.Select(bits[len(bits)-1-i], multiply, output)
	}
	return output
}
func (paillerCircuit *PaillerCircuit) Define(curveId ecc.ID, cs frontend.API) error {
	// c = g^m * r^n mod n^2 = ((m*n+1) mod n^2) * r^n mod n^2
	n := paillerCircuit.Publickey.N
	m := paillerCircuit.Data
	nsquare := paillerCircuit.Publickey.NSquared
	//((m*n+1) mod n^2)
	mn := paillerCircuit.mod(cs.Add(cs.Constant(1), cs.Mul(m, n)), nsquare, cs)
	rn := paillerCircuit.exp(cs)

	c := paillerCircuit.mod(cs.Mul(rn,
		mn), nsquare, cs)

	//fmt.Println(c.WitnessValue)
	fmt.Println(c)
	//verify
	//cs.AssertIsEqual(mn, paillerCircuit.ExpectedResult)
	return nil

}
