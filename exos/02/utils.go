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
	return hash([]byte(string(toString(structure))))
}

func hash(obj []byte) string {
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}
