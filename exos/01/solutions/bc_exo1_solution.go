package main

/**
 *-----------------------------------------
 * Exercice: 1
 *-----------------------------------------
 * Name: 	"Ledger Baker"
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

// Block defines a node in blockchain structure.
type Block struct {
	ID        int
	LastID    int // TODO previous id ...
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
	LastKey int
	Blocks  map[int]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	lastID := chain.last().ID
	b.ID = lastID + 1
	b.LastID = lastID
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
		LastKey: 0,
		Blocks:  make(map[int]Block),
	}
}

// Create a new default block with data.
func generateBlock(reference string, quantity int, price float32) Block {
	return Block{
		LastID: -1, CreatedAt: time.Now().Format(time.RFC3339),
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
