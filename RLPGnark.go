package main

import ("github.com/consensys/gnark/frontend")

func encode_item(binString){

	if isinstance(input,str):
        if len(input) == 1 and ord(input) < 0x80: return input
        else: return encode_length(len(input), 0x80) + input
    elif isinstance(input,list):
        output = ''
		for i := 0; i < len()
        for item in input: output += encode_item(item)
        return encode_length(len(output), 0xc0) + output

}

func encode_length(L,offset) {

	if L < 56:
         return chr(L + offset)
    elif L < 256**8:
         BL = to_binary(L)
         return chr(len(BL) + offset + 55) + BL
    else:
         raise Exception("input too long")

}
    

// Proposed structure
// func encodeItem(string) - it takes in a single binary string and encodes it
// func encodeList(list of strings) - takes in a list of items and encodes it
// func verifyEncoding

//encoding

// For a single byte whose value is in the [0x00, 0x7f] (decimal [0, 127]) range, that byte is its own RLP encoding.


// Otherwise, if a string is 0-55 bytes long, the RLP encoding consists of a single byte with value 0x80 (dec. 128) 
// plus the length of the string followed by the string. The range of the first byte is thus [0x80, 0xb7] (dec. [128, 183]).


// If a string is more than 55 bytes long, the RLP encoding consists of a single byte with value 0xb7 (dec. 183)
//  plus the length in bytes of the length of the string in binary form, followed by the length of the string, 
// followed by the string. For example, a 1024 byte long string would be encoded as \xb9\x04\x00 (dec. 185, 4, 0) 
// followed by the string. Here, 0xb9 (183 + 2 = 185) as the first byte, followed by the 2 bytes 0x0400 (dec. 1024) that denote the length of the actual string.
//  The range of the first byte is thus [0xb8, 0xbf] (dec. [184, 191]).


// If the total payload of a list (i.e. the combined length of all its items being RLP encoded) is 0-55 bytes long, 
// the RLP encoding consists of a single byte with value 0xc0 plus the length of the list followed by the concatenation of the RLP encodings of the items. 
// The range of the first byte is thus [0xc0, 0xf7] (dec. [192, 247]).


// If the total payload of a list is more than 55 bytes long, the RLP encoding consists of a single byte with value 0xf7 plus 
// the length in bytes of the length of the payload in binary form, followed by the length of the payload, followed by the 
// concatenation of the RLP encodings of the items. The range of the first byte is thus [0xf8, 0xff] (dec. [248, 255]).

