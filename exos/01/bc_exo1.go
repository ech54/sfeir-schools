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
	// "time" //TODO : add package
)

// Block defines a node in blockchain structure.
type Block struct {

	//TODO: finish the block structure
	// - add id, lastID, timestamp, and data

}

// Chain structure.
type Chain struct {
	lastKey int
	blocks  map[int]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {

	//TODO:
	//
}

/**
 * Enrich chain with a last block function.
 */
func (chain *Chain) last() Block {

	//TODO: review method to return the last block (current)
	// Note: use the lastKey.
	return Block{}
}

// Create the default chain structure.
func genesis() Chain {

	//TODO: review method to init the chain

	return Chain{}
}

// Create a new default block with data.
func generateBlock(price float32, quantity int, reference string) Block {

	//TODO: review method to generate a new block

	return Block{}
}

// Execute the code in console: "go run bc_exo1.go"
func main() {
	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add tree transactions:
	blockChain.addBlock(generateBlock(1.2, 5, "1111"))
	blockChain.addBlock(generateBlock(2.3, 2, "2222"))
	blockChain.addBlock(generateBlock(1.77, 4, "1111"))

	// Display blockchain
	for m, v := range blockChain.blocks {
		fmt.Println("Block: ", m, v)
	}
}
