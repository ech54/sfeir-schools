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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
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
	lastKey string
	blocks  map[string]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	b.LastID = chain.last().ID
	b.ID = hashBlock(b)
	chain.lastKey = b.ID
	chain.blocks[b.ID] = b
}

/**
 * Enrich chain with a last block function.
 */
func (chain *Chain) last() Block {
	return chain.blocks[chain.lastKey]
}

// Create the default chain structure.
func genesis() Chain {
	return Chain{
		lastKey: hash([]byte{0}),
		blocks:  make(map[string]Block),
	}
}

// Create a new default block with data.
func generateBlock(reference string, quantity int, price float32) Block {
	return Block{
		LastID: time.Now().Format(time.RFC3339), CreatedAt: time.Now().Format(time.RFC3339),
		Data: Data{Reference: reference, Quantity: quantity, Price: price},
	}
}

func hashBlock(b Block) string {
	convertedBlock, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return hash([]byte(string(convertedBlock)))
}

func hash(obj []byte) string {
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func prettyPrint(b Block) {
	j, _ := json.MarshalIndent(b, "", "  ")
	fmt.Print("\n Block: ", string(j))
}

func main() {
	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add three transactions:
	blockChain.addBlock(generateBlock("croissants", 5, 1.2))
	blockChain.addBlock(generateBlock("pains", 2, 2.3))
	blockChain.addBlock(generateBlock("croissants", 4, 1.77))

	// Display blockchain
	for _, v := range blockChain.blocks {
		prettyPrint(v)
	}
}
