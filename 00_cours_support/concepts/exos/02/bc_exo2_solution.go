package main

/**
 *-----------------------------------------
 * Exercice: 2
 *-----------------------------------------
 * Name: 	Improve block identifier
 * 			management.
 *-----------------------------------------
 * Objectives:
 *
 * - Create a simple structure to register
 * three transactions.
 * - for help on map feature, see: https://blog.golang.org/go-maps-in-action
 */
import (
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Block defines a node in blockchain structure.
type Block struct {
	ID        string
	LastID    string
	CreatedAt string
	Data      Data
}

// Data structure which contains the transaction data.
type Data struct {
	Reference string
	Quantity  int
	Price     float32
}

// Chain structure.
type Chain struct {
	LastKey string
	Blocks  map[string]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	b.LastID = chain.last().ID
	b.ID = hashStruct(b)
	chain.LastKey = b.ID
	chain.Blocks[b.ID] = b
}

/**
 * Enrich chain with a last block function.
 */
func (chain *Chain) last() Block {
	return chain.Blocks[chain.LastKey]
}

// Create the default chain structure.
func genesis() Chain {
	return Chain{
		LastKey: hash([]byte{0}),
		Blocks:  make(map[string]Block),
	}
}

// Create a new default block with data.
func generateBlock(reference string, quantity int, price float32) Block {
	return Block{
		LastID: time.Now().Format(time.RFC3339), CreatedAt: time.Now().Format(time.RFC3339),
		Data: Data{Reference: reference, Quantity: quantity, Price: price},
	}
}

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

func hashStructBytes(structure interface{}) []byte {
	return hashBytes([]byte(string(toString(structure))))
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

func main() {
	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add three transactions:
	blockChain.addBlock(generateBlock("croissants", 5, 1.2))
	blockChain.addBlock(generateBlock("pains", 2, 2.3))
	blockChain.addBlock(generateBlock("croissants", 4, 1.77))

	// Display blockchain
	for _, v := range blockChain.Blocks {
		prettyPrint(v)
	}
}
