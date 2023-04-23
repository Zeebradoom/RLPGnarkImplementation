import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)


type DecodingCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"encoded"`
	Y frontend.Variable `gnark:"decoded"`
}

func (circuit *DecodingCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	return nil
}


// type Circuit interface {
//     // Define declares the circuit's Constraints
//     Define(api frontend.API) error
// }

