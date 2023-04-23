package main

import (
    "fmt"
    "math/big"

    "github.com/consensys/gnark/backend/groth16"
    "github.com/consensys/gnark/frontend"
    "github.com/ethereum/go-ethereum/rlp"
)

type RLP []byte

func (r RLP) Decode() (interface{}, error) {
    if len(r) == 0 {
        return []interface{}{}, nil
    }
    if r[0] <= 0x7f {
        return r[0], nil
    }
    if r[0] <= 0xb7 {
        return r[1 : 1+r[0]-0x80], nil
    }
    if r[0] <= 0xbf {
        l, _, err := rlp.DecodeSize(r)
        if err != nil {
            return nil, err
        }
        return r[1+l:], nil
    }
    if r[0] <= 0xf7 {
        l := r[0] - 0xc0
        res := make([]interface{}, l)
        for i := 0; i < int(l); i++ {
            item, err := RLP(r[1:]).Decode()
            if err != nil {
                return nil, err
            }
            res[i] = item
            r = r[1+item.([]byte)[0]]
        }
        return res, nil
    }
    l, _, err := rlp.DecodeSize(r)
    if err != nil {
        return nil, err
    }
    res := make([]interface{}, l)
    for i := 0; i < int(l); i++ {
        item, err := RLP(r[1:]).Decode()
        if err != nil {
            return nil, err
        }
        res[i] = item
        r = r[1+item.(RLP)[0]]
    }
    return res, nil
}

func main() {
    // Encode a nested array
    data := []interface{}{1, []interface{}{2, 3}, 4}
    encoded, err := rlp.EncodeToBytes(data)
    if err != nil {
        panic(err)
    }

    // Define the circuit
    circuit := frontend.New()

    // Define the inputs
    encodedVar := circuit.PUBLIC_INPUT("encoded")
    decodedVar := circuit.SECRET_INPUT("decoded")

    // Decode the RLP input
    decoded := new(big.Int)
    err = circuit.Define(decoded, func() error {
        decodedData, err := RLP(encodedVar.Get()).Decode()
        if err != nil {
            return err
        }
        return decodedVar.Assign(decodedData)
    })
    if err != nil {
        panic(err)
    }

    // Compile and generate the proving and verification keys
    pk, vk, err := groth16.Setup(circuit)
    if err != nil {
        panic(err)
    }

    // Generate a proof
    proof, err := groth16.Prove(circuit, pk, encodedVar, decodedVar)
    if err != nil {
        panic(err)
    }

    // Verify the proof
    result, err := groth16.Verify(proof, vk, encodedVar, decodedVar)
	if err != nil {
		panic(err)
	}
	if !result {
		panic("proof is invalid")
	}
	fmt.Println("RLP decoding succeeded:", decoded.Var)
}