package zktx

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

type ModCircuit struct {
	X              frontend.Variable
	Y              frontend.Variable
	ExpectedResult frontend.Variable
}

func (ModCircuit *ModCircuit) Define(curveId ecc.ID, api frontend.API) error {
	v1 := ModCircuit.X
	v2 := ModCircuit.Y
	v3 := api.Mod(v1, v2)
	var x = frontend.Value(v3)

	api.AssertIsEqual(x, ModCircuit.ExpectedResult)
	return nil
}
