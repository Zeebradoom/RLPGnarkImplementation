package main

import (
	"github.com/consensys/gnark/frontend"
)

func ArrayEq(api frontend.API, a, b []frontend.Variable, inLen frontend.Variable, nIn int) (frontend.Variable, error) {
	api.AssertIsLessOrEqual(inLen, 252)
	api.AssertIsLessOrEqual(nIn, 252) //is 252 in integer or bits?
	
	matchSum := make([]frontend.Variable, nIn)

	for idx := 0; idx < nIn; idx++ {
		i := api.Sub(a[idx], b[idx])

		if idx == 0 {
			matchSum[idx] = i
		} else {
			api.Add(matchSum[idx], matchSum[idx-1], i)
		}
	}

	matchChooser:= Multiplexer( matchSum, inLen)
	// matchChooser[0][0] := 0 //not sure, theres no documentation

	for idx := 0; idx < nIn; idx++ {
		matchChooser[idx+1][0]:= matchSum[idx]
	}
	// matchChooser.Sel := inLen //idk

	matchCheck := api.sub(matchChooser[0], inLen)
	return api.IsZero(matchCheck)
}

func EscalarProduct(api frontend.API, in1, in2 []frontend.Variable) frontend.Variable {
	out := frontend.Variable(0)
	for i := 0; i < len(in1); i++ {
		out = api.Add(out, api.Mul(in1[i], in2[i]))
	}
	return out
}

func Decoder(api frontend.API, inp frontend.Variable, w int) ([]frontend.Variable, frontend.Variable) {
	out := make([]frontend.Variable, w)
	success := frontend.Variable(0)
	for i := 0; i < w; i++ {
		out[i] = IsEqual(api, i, inp)
		api.AssertIsEqual(api.Mul(out[i], api.Sub(out[i], 1)), 0)
		success = api.Add(success, out[i])
	}
	api.AssertIsEqual(api.Mul(success, api.Sub(success, 1)), 0)
	return out, success
}

func Multiplexer(api frontend.API, inp [][]frontend.Variable, sel frontend.Variable) []frontend.Variable {
	nIn := len(inp)
	wIn := len(inp[0])
	out := make([]frontend.Variable, wIn)

	decodedData, success := Decoder(api, sel, nIn)
	api.AssertIsEqual(success, 1)

	for j := 0; j < wIn; j++ {
		in1 := make([]frontend.Variable, nIn)
		//in2 := make([]frontend.Variable, nIn)
		for k := 0; k < nIn; k++ {
			in1[k] = inp[k][j]
			//in2[k] = decodedData[k]
		}
		out[j] = EscalarProduct(api, in1, decodedData)
	}
	return out
}