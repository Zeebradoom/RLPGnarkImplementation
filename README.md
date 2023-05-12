**RLP Encoding using Gnark**
process design and implementation - used Ethereum's description found here - https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/#:~:text=RLP%20standardizes%20the%20transfer%20of,objects%20in%20Ethereum's%20execution%20layer

The aim of your work is to decode the compressed transaction within the circuit, before you can do other computation on it also within the circuit.

Inspired by Circom circuit implementation by the  Yi-Sun, found here: https://github.com/yi-sun/zk-attestor 