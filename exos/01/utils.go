package main

import (
	"encoding/json"
	"fmt"
)

func prettyPrint(b interface{}) {
	j, _ := json.MarshalIndent(b, "", "  ")
	fmt.Print("\n Block: ", string(j))
}
