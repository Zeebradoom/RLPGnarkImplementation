package main

import (
	//"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/liyue201/gnark-circomlib/blob/main/circuits/multiplexer.go"
)

type EncodingCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	// check decoded = RLP-decode(encoded))
	X frontend.Variable `gnark:"encoded"`
	Y frontend.Variable `gnark:"decoded"`
} 

func SubArray(api frontend.API, nIn, maxSelect, nInBits int) (out []frontend.Variable, outLen frontend.Variable) {

	// define the circuit
	circuit := frontend.New();

	// define inputs
	in := make([]frontend.Variable, nIn)
	for i := 0; i < nIn; i++ {
		in[i] = circuit.PbInput() // public input
	}
	start := circuit.PbInput() // public input
	end := circuit.PbInput()   // public input

	// define outputs
	out = make([]frontend.Variable, maxSelect)
	for i := 0; i < maxSelect; i++ {
		out[i] = circuit.PbOutput() // public output
	}
	outLen = circuit.PbOutput() // public output

	// define components
	lt1 := circuit.LessEq(start, end)                    // start <= end
	lt2 := circuit.LessEq(end, circuit.Uint64(nIn-1))     // end <= nIn-1
	lt3 := circuit.LessEq(circuit.Sub(end, start), maxSelect) // end - start <= maxSelect
	n2b := circuit.ToBinary(start, nInBits)              // convert start to a bit array
	shifts := make([][]frontend.Variable, nInBits)        // shifts variable for each bit position
	for i := 0; i < nInBits; i++ {
		shifts[i] = make([]frontend.Variable, nIn)
		for j := 0; j < nIn; j++ {
			if i == 0 {
				tempIdx := (j + (1 << i)) % nIn
				shifts[i][j] = circuit.Mul(n2b[i], circuit.Sub(in[tempIdx], in[j]))
				shifts[i][j] = circuit.Add(shifts[i][j], in[j])
			} else {
				prevIdx := i - 1
				tempIdx := (j + (1 << i)) % nIn
				shifts[i][j] = circuit.Mul(n2b[i], circuit.Sub(shifts[prevIdx][tempIdx], shifts[prevIdx][j]))
				shifts[i][j] = circuit.Add(shifts[i][j], shifts[prevIdx][j])
			}
		}
	}

	// define constraints
	circuit.Assert(lt1)                         // start <= end
	circuit.Assert(lt2)                         // end <= nIn-1
	circuit.Assert(lt3)                         // end - start <= maxSelect
	circuit.AssertIsEqual(outLen, circuit.Sub(end, start)) // outLen = end - start
	for i := 0; i < maxSelect; i++ {
		circuit.AssertIsEqual(out[i], shifts[nInBits-1][i]) // out[i] = shifts[nInBits-1][i]
	}

	return out, outLen
}

func ShiftRight(circuit *frontend.Circuit, nIn int, nInBits uint64) Define(api frontend.API){
    in := make([]frontend.Signal, nIn)
    for i := 0; i < nIn; i++ {
        in[i] = circuit.SIGNAL(fmt.Sprintf("in_%v", i))
    }

    shift := circuit.SIGNAL("shift")
    out := make([]frontend.Signal, nIn)
    for i := 0; i < nIn; i++ {
        out[i] = circuit.SIGNAL(fmt.Sprintf("out_%v", i))
    }

    n2b := frontend.ToBinary(circuit, nInBits, shift)

    shifts := make([][]frontend.Signal, nInBits)
    for i := range shifts {
        shifts[i] = make([]frontend.Signal, nIn)
    }

    for idx := 0; idx < int(nInBits); idx++ {
        if idx == 0 {
            for j := 0; j < min((1 << idx), nIn); j++ {
                circuit.MUL(n2b[idx], circuit.SUB(circuit.ZERO, in[j]), circuit.ONE).ADD(in[j], shifts[0][j])
            }
            for j := (1 << idx); j < nIn; j++ {
                tempIdx := j - (1 << idx)
                circuit.MUL(n2b[idx], circuit.SUB(in[tempIdx], in[j]), circuit.ONE).ADD(in[j], shifts[0][j])
            }
        } else {
            prevIdx := idx - 1
            for j := 0; j < min((1 << idx), nIn); j++ {
                circuit.MUL(n2b[idx], circuit.SUB(circuit.ZERO, shifts[prevIdx][j]), circuit.ONE).ADD(shifts[prevIdx][j], shifts[idx][j])
            }
            for j := (1 << idx); j < nIn; j++ {
                tempIdx := j - (1 << idx)
                circuit.MUL(n2b[idx], circuit.SUB(shifts[prevIdx][tempIdx], shifts[prevIdx][j]), circuit.ONE).ADD(shifts[prevIdx][j], shifts[idx][j])
            }
        }
    }

    for i := 0; i < nIn; i++ {
        circuit.TO_BINARY(out[i], shifts[nInBits-1][i])
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// func (circuit *EncodingCircuit) Define(api frontend.API) error {


// 	return nil;
// }

func (circuit *DecodingCircuit) ArrayEq(api frontend.API) error {

}

func ArrayEq(api frontend.API,circuit multiplexer, a, b []frontend.Variable, inLen frontend.Variable, nIn int) error {
	api.AssertIsLessOrEqual(inLen, 252)
	api.AssertIsLessOrEqual(nIn, 252) //is 252 in integer or bits?
	
	matchSum := make([]frontend.Variable, nIn)

	for idx := 0; idx < nIn; idx++ {
		i := api.SUB(a[idx]-b[idx])

		if idx == 0 {
			matchSum[idx].Equal(i)
		} else {
			api.ADD(matchSum[idx], matchSum[idx-1], i)
		}
	}

	matchChooser := circuit.Multiplexer(1, nIn+1)
	matchChooser[0][0] := 0 //not sure, theres no documentation

	for idx := 0; idx < nIn; idx++ {
		matchChooser[idx+1][0]:= matchSum[idx]
	}
	// matchChooser.Sel := inLen //idk

	matchCheck := api.sub(matchChooser[0], inLen)
	return api.IsZero(matchCheck)
}

















