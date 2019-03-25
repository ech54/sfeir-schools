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
)

// Block structure.
type Block struct {
	ID        string // [1] int -> string
	LastID    string // [2] int -> string
	CreatedAt string
	Data      Data
}

// Data structure.
type Data struct {
	Reference string
	Quantity  int
	Price     float32
}

// Chain structure.
type Chain struct {
	LastKey string
	Blocks  map[string]Block // [3] int -> string
}

// add a new block to chain.
func (chain *Chain) addBlock(b Block) {
	b.LastID = chain.last().ID
	b.ID = "" // -> TODO use : hashStruct(b)
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
		LastKey: "", // -> TODO use: hash([]byte{0})
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

func main() {
	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add three transactions:
	blockChain.addBlock(generateBlock("croissants", 5, 1.2))
	blockChain.addBlock(generateBlock("pains", 2, 2.3))
	blockChain.addBlock(generateBlock("croissants", 4, 1.77))

	// Display blockchain
	for _, v := range blockChain.Blocks {
		fmt.Println("Block: ", v)
	}
}
