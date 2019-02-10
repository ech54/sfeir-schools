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
	id        int
	lastID    int // TODO previous id ...
	createdAt time.Time
	data      Data
}

// Data structure which contains the transaction data.
type Data struct {
	quantity  int
	reference string
	price     float32
}

// Chain structure.
type Chain struct {
	lastKey int
	blocks  map[int]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	lastID := chain.last().id
	b.id = lastID + 1
	b.lastID = lastID
	chain.lastKey = b.id
	chain.blocks[b.id] = b
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
		lastKey: 0,
		blocks:  make(map[int]Block),
	}
}

// Create a new default block with data.
func generateBlock(price float32, quantity int, reference string) Block {
	return Block{
		lastID: -1, createdAt: time.Now(),
		data: Data{price: price, quantity: quantity, reference: reference},
	}
}

func main() {
	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add tree transactions:
	blockChain.addBlock(generateBlock(1.2, 5, "croissants"))
	blockChain.addBlock(generateBlock(2.3, 2, "pains"))
	blockChain.addBlock(generateBlock(1.77, 4, "croissants"))

	// Display blockchain
	for m, v := range blockChain.blocks {
		fmt.Println("Block: ", m, v)
	}
}
