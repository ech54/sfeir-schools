package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func prettyPrint(structure interface{}) {
	j, _ := json.MarshalIndent(structure, "", "  ")
	fmt.Print("\n Block: ", string(j))
}

func toString(structure interface{}) string {
	bytesStructure, err := json.Marshal(structure)
	if err != nil {
		panic(err)
	}
	return string(bytesStructure)
}

func hashStruct(structure interface{}) string {
	return hash([]byte(toString(structure)))
}

func HashStructBytes(structure interface{}) ([]byte, error) {
	return hashBytes([]byte(toString(structure))), nil
	//return []byte{0}
}

func hash(obj []byte) string {
	return hex.EncodeToString(hashBytes(obj))
}

func hashBytes(obj []byte) []byte {
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
