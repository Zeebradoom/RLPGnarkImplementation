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
	//get most significant, hex bit that tells us the length of the bits
	bits := api.ToBinary(circuit.X, 8)
	prefix := api.FromBinary(bits[:]...)

	bound_bits := api.ToBinary(0x7f)
	bound := api.FromBinary(bound_bits[:]...)

	api.AssertIsLessOrEqual(prefix, bound)
	api.AssertIsEqual(circuit.Y, prefix)
	// if cond {
		// return api.AssertIsEqual(circuit.Y, prefix)
	// } else{
		return nil
	// }
}

func main() {
	// compiles our circuit into a R1CS
	var circuit DecodingCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// // groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// // witness definition
	assignment := DecodingCircuit{X: 0x7f, Y: 0x7f}
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

// if prefix <= 0x7f {
// 	return string([]byte{prefix}), index, nil
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