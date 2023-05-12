package main

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type DecodingCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	// check decoded = RLP-decode(encoded))
	X frontend.Variable `gnark:"encoded"`
	Y frontend.Variable `gnark:"decoded"`
} 

func (circuit *DecodingCircuit) Define(api frontend.API) error {
	//get the prefix
	bits := api.ToBinary(circuit.X, 8)
	prefix := api.FromBinary(bits[:]...)

	// bound is a bound for a [0:127]
	bound_bits := api.ToBinary(0x7f)
	bound := api.FromBinary(bound_bits[:]...)
	
	compare := api.Cmp(prefix, bound)
	//fmt.Printf("Compare: %v\n", compare)
	is_equal := api.IsZero(compare)

	one_bits := api.ToBinary(0x01)
	one := api.FromBinary(one_bits[:]...)
	
	diff := api.Sub(compare, one)
	is_less := api.IsZero(diff)

	check := api.Or(is_equal, is_less)
	//check

	api.AssertIsEqual(check, one)
	api.AssertIsLessOrEqual(prefix, bound)
	api.AssertIsEqual(circuit.Y, prefix)

	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit DecodingCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// // groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// // witness definition
	assignment := DecodingCircuit{X: 0x07, Y: 0x07}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	fmt.Printf("Witness: %v\n", witness)
	publicWitness, _ := witness.Public()

	// // groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)

	fmt.Printf("Main works")

}


// prefix := input[index]
// index++

// } else if prefix <= 0xb7 {
// 	length := int(prefix - 0x80)
// 	if index+length > len(input) {
// 		return nil, index, fmt.Errorf("String length out of range")
// 	}

// 	return string(input[index : index+length]), index + length, nil
// } else if prefix <= 0xbf {
// 	lengthLength := int(prefix - 0xb7)
// 	if index+lengthLength > len(input) {
// 		return nil, index, fmt.Errorf("Length out of range")
// 	}

// 	length := bytesToInt(input[index : index+lengthLength])
// 	index += lengthLength
// 	if index+length > len(input) {
// 		return nil, index, fmt.Errorf("String length out of range")
// 	}

// 	return string(input[index : index+length]), index + length, nil




// if_same := api.cmp()
	//bound_less_than_183_bits := api.ToBinary(0xb7)
	//bound_less_than_183 := api.FromBinary(bound_less_than_183_bits[:]...)
