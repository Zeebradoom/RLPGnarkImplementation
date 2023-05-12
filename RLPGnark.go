package main

import (
    "fmt"
"strings"
)

// type Circuit struct {
// 	Secret frontend.Variable
// 	Hash   frontend.Variable `gnark:",public"` // struct tags default visibility is "secret"
// }

func rlpSerialize(item interface{}) string {
	switch v := item.(type) {
	case string:
		return serializeString(v)
	case []interface{}:
		return serializeList(v)
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func serializeString(input string) string {
	inputBytes := []byte(input)
	inputLength := len(inputBytes)

	if inputLength == 1 && inputBytes[0] <= 0x7f {
		return input
	}

	var result strings.Builder

	if inputLength <= 55 {
		result.WriteByte(byte(0x80 + inputLength))
	} else {
		lengthBytes := intToBytes(inputLength)
		lengthBytesLength := len(lengthBytes)

		result.WriteByte(byte(0xb7 + lengthBytesLength))
		result.WriteString(string(lengthBytes))
	}

	result.WriteString(input)
	return result.String()
}

func serializeList(items []interface{}) string {
	serializedItems := make([]string, len(items))
	for i, item := range items {
		serializedItems[i] = rlpSerialize(item)
	}

	serializedList := strings.Join(serializedItems, "")
	listLength := len(serializedList)

	var result strings.Builder

	if listLength <= 55 {
		result.WriteByte(byte(0xc0 + listLength))
	} else {
		lengthBytes := intToBytes(listLength)
		lengthBytesLength := len(lengthBytes)

		result.WriteByte(byte(0xf7 + lengthBytesLength))
		result.WriteString(string(lengthBytes))
	}

	result.WriteString(serializedList)
	return result.String()
}

func intToBytes(n int) []byte {
	var result []byte
	for n > 0 {
		result = append([]byte{byte(n % 256)}, result...)
		n /= 256
	}
	return result
}

func rlpDecode(input string) (interface{}, error) {
	decoded, _, err := decode([]byte(input), 0)
	return decoded, err
}

func decode(input []byte, index int) (interface{}, int, error) {
	if index >= len(input) {
		return nil, index, fmt.Errorf("Unexpected end of input")
	}

	prefix := input[index]
	index++

	if prefix <= 0x7f {
		return string([]byte{prefix}), index, nil
	} else if prefix <= 0xb7 {
		length := int(prefix - 0x80)
		if index+length > len(input) {
			return nil, index, fmt.Errorf("String length out of range")
		}

		return string(input[index : index+length]), index + length, nil
	} else if prefix <= 0xbf {
		lengthLength := int(prefix - 0xb7)
		if index+lengthLength > len(input) {
			return nil, index, fmt.Errorf("Length out of range")
		}

		length := bytesToInt(input[index : index+lengthLength])
		index += lengthLength
		if index+length > len(input) {
			return nil, index, fmt.Errorf("String length out of range")
		}

		return string(input[index : index+length]), index + length, nil
	} else if prefix <= 0xf7 {
		length := int(prefix - 0xc0)
		return decodeList(input, index, index+length)
	} else {
		lengthLength := int(prefix - 0xf7)
		if index+lengthLength > len(input) {
			return nil, index, fmt.Errorf("Length out of range")
		}

		length := bytesToInt(input[index : index+lengthLength])
		index += lengthLength
		return decodeList(input, index, index+length)
	}
}

func decodeList(input []byte, index, end int) ([]interface{}, int, error) {
	items := make([]interface{}, 0)

	for index < end {
		item, newIndex, err := decode(input, index)
		if err != nil {
			return nil, newIndex, err
		}
		items = append(items, item)
		index = newIndex
	}

	return items, index, nil
}

func bytesToInt(bytes []byte) int {
	result := 0
	for _, b := range bytes {
		result = (result << 8) + int(b)
	}
	return result
}

func main() {
	items := []interface{}{
		"cat",
		[]interface{}{"puppy", "cow"},
		"horse",
		[]interface{}{[]interface{}{}},
		"pig",
		[]interface{}{""},
		"sheep",
	}

	encoded := rlpSerialize(items)
    decoded, err := rlpDecode(encoded)
    if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Items: %v\nEncoded: %x\n Decoded: %v\n", items, encoded, decoded)
	}

    item := "hello, world"

	encoded_single := rlpSerialize(item)
    decoded_single, err := rlpDecode(encoded_single)
    if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Items: %v\nEncoded: %x\n Decoded: %v\n", item, encoded_single, decoded_single)
	}
	
}

// func rlpDecodeUniversal(input1 string, inout2 string)

func rlpCheckDecode(input1 string, input2 string) {
    // check input2 = RLP-decode(input1))
    
}
