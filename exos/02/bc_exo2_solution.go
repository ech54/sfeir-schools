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
	lastID    string
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
	lastKey string
	blocks  map[string]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	lastID := chain.last().ID
	b.ID = hashBlock(b)
	b.lastID = lastID
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
		lastKey: "",
		blocks:  make(map[string]Block),
	}
}

// Create a new default block with data.
func generateBlock(price float32, quantity int, reference string) Block {
	return Block{
		lastID: "", createdAt: time.Now(),
		data: Data{price: price, quantity: quantity, reference: reference},
	}
}

func hashBlock(b Block) string {
	convertedBlock, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(convertedBlock))
	return hash([]byte(string(convertedBlock)))
}

func hash(obj []byte) string {
	fmt.Println(obj)
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
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
