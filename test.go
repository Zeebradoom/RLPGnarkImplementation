package main

import (
	"fmt"

	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/consensys/gnark-crypto/ecc"
)

type RLPDecodingCircuit struct {
	Input  frontend.Variable
	Output frontend.Variable
}

func (circuit *RLPDecodingCircuit) Define(api frontend.API) error {
	// unpack the input variable to a slice of bytes
	inputBytes := cs.ToBytes(circuit.Input)

	// decode the input using RLP
	var decoded []interface{}
	err := rlp.DecodeBytes(inputBytes, &decoded)
	if err != nil {
		return err
	}

	// pack the decoded data into the output variable
	outputBytes := make([]byte, 0)
	for _, item := range decoded {
		switch value := item.(type) {
		case []byte:
			outputBytes = append(outputBytes, value...)
		default:
			return fmt.Errorf("unexpected type in RLP-decoded data: %T", value)
		}
	}
	circuit.Output.Assign(cs, backend.UInt64LittleEndian(outputBytes))

	return nil
}

func DecodeBytes(encoded []byte) ([]byte, error) {
	// create a new circuit instance
	circuit := &RLPDecodingCircuit{}

	// create a new R1CS instance from the circuit definition
	r1cs, err := frontend.Compile(groth16.NewCurveParams().ID, circuit)
	if err != nil {
		return nil, err
	}

	// create a new prover instance
	prover := groth16.NewProver()

	// create a new input variable and assign the input data
	inputVar := circuit.Input
	inputVar.AssignBytes(r1cs.ToFieldElement(encoded))

	// generate a proof for the input data
	proof, err := prover.GenerateProof(r1cs, []*frontend.ConstraintSystem{circuit}, groth16.DefaultProvingKey)
	if err != nil {
		return nil, err
	}

	// verify the proof
	err = groth16.Verify(proof, r1cs, groth16.DefaultPublicInputs(encoded))
	if err != nil {
		return nil, err
	}

	// extract the output from the circuit instance
	outputBytes := circuit.Output.GetBackend().(backend.UInt64LittleEndian)

	return outputBytes, nil
}

func main() {
	// create a sample RLP-encoded object
	rlpEncoded := []byte{0xc8, 0x86, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa2, 0x06, 0x00}

	// decode the RLP-encoded object using gnark
	decoded, err := DecodeBytes(rlpEncoded)
	if err != nil {
		panic(err)
	}

	// print the decoded data
	fmt.Printf("%x\n", decoded)
}